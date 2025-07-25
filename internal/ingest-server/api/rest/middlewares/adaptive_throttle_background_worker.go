package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/ingest-server/contract"
)

// RPSStatsWorker is a background worker that periodically fetches throttle status,
// muted fingerprints, and calculates RPS from Redis.
type RPSStatsWorker struct {
	redisClient  *redis.Client
	cache        *RPSStatsCache
	projectsRepo contract.ProjectsRepository
	interval     time.Duration
	rpsWindow    time.Duration
}

// NewRPSStatsWorker creates a new RPSStatsWorker.
func NewRPSStatsWorker(
	redisClient *redis.Client,
	cache *RPSStatsCache,
	projectsRepo contract.ProjectsRepository,
	interval time.Duration,
	rpsWindow time.Duration,
) *RPSStatsWorker {
	return &RPSStatsWorker{
		redisClient:  redisClient,
		cache:        cache,
		projectsRepo: projectsRepo,
		interval:     interval,
		rpsWindow:    rpsWindow,
	}
}

// Start starts the worker.
func (w *RPSStatsWorker) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.fetchAllProjectThrottleStatus(ctx)
			}
		}
	}()
}

// fetchAllProjectThrottleStatus fetches throttle status and calculates RPS for all projects from Redis.
func (w *RPSStatsWorker) fetchAllProjectThrottleStatus(ctx context.Context) {
	// Get a list of all projects from the repository
	projectIDs, err := w.projectsRepo.GetProjectIDs(ctx)
	if err != nil {
		slog.Error("Error getting project IDs from repository", "error", err)

		// Fallback to using the cache if the repository call fails
		w.cache.mu.RLock()
		// Collect project IDs from both maps
		projectIDsMap := make(map[domain.ProjectID]struct{})
		for projectID := range w.cache.rpsStats {
			projectIDsMap[projectID] = struct{}{}
		}
		for projectID := range w.cache.throttleEnabled {
			projectIDsMap[projectID] = struct{}{}
		}

		// Convert to slice
		projectIDs = make([]domain.ProjectID, 0, len(projectIDsMap))
		for projectID := range projectIDsMap {
			projectIDs = append(projectIDs, projectID)
		}
		w.cache.mu.RUnlock()
	}

	// Calculate RPS for each project in parallel
	var wg sync.WaitGroup
	for _, projectID := range projectIDs {
		wg.Add(1)
		go func(pid domain.ProjectID) {
			defer wg.Done()

			w.calculateProjectRPS(ctx, pid)
		}(projectID)
	}
	wg.Wait()

	// Log current RPS values
	// w.logRPSStats()
}

// calculateProjectRPS calculates the RPS for a project.
func (w *RPSStatsWorker) calculateProjectRPS(ctx context.Context, projectID domain.ProjectID) {
	// Calculate the current time batch
	now := time.Now()
	currentBatch := calculateTimeBatch(now, w.rpsWindow)

	// Check if the batch has changed
	storedBatch, ok := w.cache.GetCurrentBatch(projectID)
	if ok && storedBatch != currentBatch {
		// Batch has changed, store the previous RPS value
		prevCnt, err := w.redisClient.Get(ctx, generateRPSBatchKey(projectID, storedBatch)).Uint64()
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				slog.Error("Error getting RPS counter", "error", err)
			}

			prevCnt = 0
		}

		prevRPS := CalculateRPS(prevCnt, w.rpsWindow)
		w.cache.SetPrevBatchRPS(projectID, prevRPS)
	}

	// Update the current batch
	w.cache.SetCurrentBatch(projectID, currentBatch)

	// Get the counter from Redis
	key := generateRPSBatchKey(projectID, currentBatch)
	count, err := w.redisClient.Get(ctx, key).Uint64()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			slog.Error("Error getting RPS counter", "error", err)

			return
		}

		// If the key doesn't exist, set the count to 0
		count = 0
	}

	// Calculate the RPS
	rps := CalculateRPS(count, w.rpsWindow)

	if w.rpsWindow.Seconds() > 1 {
		prevRPS, ok := w.cache.GetPrevBatchRPS(projectID)
		if ok {
			rpsWindowSec := w.rpsWindow.Seconds()
			secondsLeftInBatch := timeLeftInBatch(currentBatch, w.rpsWindow).Seconds()
			rps += uint64(float64(prevRPS) * (rpsWindowSec - secondsLeftInBatch) / rpsWindowSec * influenceCoeff)
		}
	}

	// Update the cache
	w.cache.SetRPS(projectID, rps)
}

// logRPSStats logs the current RPS values for all projects.
//func (w *RPSStatsWorker) logRPSStats() {
//	w.cache.mu.RLock()
//	defer w.cache.mu.RUnlock()
//
//	if len(w.cache.rpsStats) == 0 {
//		slog.Info("No RPS stats available")
//
//		return
//	}
//
//	data := make([][2]string, 0, len(w.cache.rpsStats))
//	for projectID, rps := range w.cache.rpsStats {
//		data = append(data, [2]string{
//			strconv.FormatUint(uint64(projectID), 10),
//			strconv.FormatUint(rps, 10),
//		})
//	}
//
//	sort.Slice(data, func(i, j int) bool {
//		return data[i][0] < data[j][0]
//	})
//
//	slog.Debug("Current RPS", "rps", data)
//}

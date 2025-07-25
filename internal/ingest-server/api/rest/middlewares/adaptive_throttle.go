package middlewares

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/ingest-server/contract"
)

const influenceCoeff = 0.95

// incrRPSScript is a LUA script that increments a counter and sets its TTL if it's the first increment.
var incrRPSScript = redis.NewScript(`
	local current = redis.call("INCR", KEYS[1])
	if current == 1 then
		redis.call("EXPIRE", KEYS[1], tonumber(ARGV[1]))
	end
	return current
`)

// AdaptiveThrottleConfig holds configuration for the adaptive throttle middleware.
type AdaptiveThrottleConfig struct {
	RPSWindow            time.Duration
	StatsRefreshInterval time.Duration
	RateLimit            uint64
}

// DefaultAdaptiveThrottleConfig returns the default configuration for the adaptive throttle middleware.
func DefaultAdaptiveThrottleConfig() *AdaptiveThrottleConfig {
	return &AdaptiveThrottleConfig{
		RPSWindow:            time.Minute,
		StatsRefreshInterval: time.Second,
		RateLimit:            100,
	}
}

// RPSStats holds the RPS and rate limit values for a project.
type RPSStats struct {
	RPS             uint64
	RateLimit       uint64
	ThrottleEnabled bool
}

// AdaptiveThrottleResult contains the middleware handler and the worker instance.
type AdaptiveThrottleResult struct {
	Handler func(http.Handler) http.Handler
	Worker  *RPSStatsWorker
}

// AdaptiveThrottle is a middleware that implements adaptive throttling for envelope endpoints.
// It returns both the middleware handler and the worker instance so
// that the application can stop it when shutting down.
func AdaptiveThrottle(
	ctx context.Context,
	redisClient *redis.Client,
	projectsRepo contract.ProjectsRepository,
	config *AdaptiveThrottleConfig,
) AdaptiveThrottleResult {
	if config == nil {
		config = DefaultAdaptiveThrottleConfig()
	}

	// Load the increment RPS script into Redis
	if err := incrRPSScript.Load(ctx, redisClient).Err(); err != nil {
		slog.Error("Failed to load increment RPS script", "error", err)
	}

	// Create the RPS stats cache and muted fingerprints cache
	cache := NewRPSStatsCache(config.RateLimit)

	// Create and start the worker
	worker := NewRPSStatsWorker(
		redisClient,
		cache,
		projectsRepo,
		config.StatsRefreshInterval,
		config.RPSWindow,
	)
	worker.Start(ctx)

	// Create the middleware handler
	handler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			projectID := wardencontext.ProjectID(ctx)

			// Check if we're over the RPS limit
			if isOverRPS(projectID, cache) {
				respond429(writer, "Global rate limit hit")

				return
			}

			// Increment the counter in Redis
			err := incrementRPSCounter(ctx, redisClient, projectID, config.RPSWindow)
			if err != nil {
				slog.Error("Error incrementing RPS counter", "error", err)
			}

			next.ServeHTTP(writer, req)
		})
	}

	return AdaptiveThrottleResult{
		Handler: handler,
		Worker:  worker,
	}
}

// isOverRPS checks if the project is over the RPS limit.
func isOverRPS(
	projectID domain.ProjectID,
	cache *RPSStatsCache,
) bool {
	// Get the RPS from the cache (calculated by the background worker)
	rps, ok := cache.GetRPS(projectID)
	if !ok {
		// If RPS is not in the cache, assume it's not over the limit
		return false
	}

	return rps > cache.rateLimit
}

// respond429 responds with a 429 Too Many Requests status code.
func respond429(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	resp := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Error encoding response", "error", err)
	}
}

// generateRPSBatchKey creates a Redis key for project RPS stats with time batch.
func generateRPSBatchKey(projectID domain.ProjectID, timeBatch int64) string {
	return "stats:project:" + strconv.FormatUint(uint64(projectID), 10) + ":rps:" +
		strconv.FormatInt(timeBatch, 10)
}

// calculateTimeBatch calculates the time batch based on the current time and the RPSWindow.
func calculateTimeBatch(now time.Time, rpsWindow time.Duration) int64 {
	return now.Unix() / int64(rpsWindow.Seconds())
}

func batchToTime(timeBatch int64, rpsWindow time.Duration) time.Time {
	return time.Unix(timeBatch*int64(rpsWindow.Seconds()), 0)
}

func timeLeftInBatch(timeBatch int64, rpsWindow time.Duration) time.Duration {
	return time.Since(batchToTime(timeBatch, rpsWindow))
}

// incrementRPSCounter increments the RPS counter for a project and sets the TTL.
func incrementRPSCounter(
	ctx context.Context,
	redisClient *redis.Client,
	projectID domain.ProjectID,
	rpsWindow time.Duration,
) error {
	// Calculate the time batch
	currentBatch := calculateTimeBatch(time.Now(), rpsWindow)

	// Generate the key
	key := generateRPSBatchKey(projectID, currentBatch)

	// Run the script
	return incrRPSScript.Run(ctx, redisClient, []string{key}, int(rpsWindow.Seconds()*2)).Err()
}

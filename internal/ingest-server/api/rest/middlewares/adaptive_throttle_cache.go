package middlewares

import (
	"sync"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

// RPSStatsCache is a cache for RPS stats.
type RPSStatsCache struct {
	mu              sync.RWMutex
	rpsStats        map[domain.ProjectID]uint64
	throttleEnabled map[domain.ProjectID]bool
	rateLimit       uint64
	// Store previous batch RPS values
	prevBatchRPS map[domain.ProjectID]uint64
	// Store the current time batch for each project
	currentBatch map[domain.ProjectID]int64
}

// NewRPSStatsCache creates a new RPSStatsCache.
func NewRPSStatsCache(rateLimit uint64) *RPSStatsCache {
	return &RPSStatsCache{
		rpsStats:        make(map[domain.ProjectID]uint64),
		throttleEnabled: make(map[domain.ProjectID]bool),
		rateLimit:       rateLimit,
		prevBatchRPS:    make(map[domain.ProjectID]uint64),
		currentBatch:    make(map[domain.ProjectID]int64),
	}
}

// GetRPS returns the RPS for a project.
func (c *RPSStatsCache) GetRPS(projectID domain.ProjectID) (uint64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	rps, ok := c.rpsStats[projectID]

	return rps, ok
}

// SetRPS sets the RPS for a project.
func (c *RPSStatsCache) SetRPS(projectID domain.ProjectID, rps uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rpsStats[projectID] = rps
}

// SetPrevBatchRPS sets the previous batch RPS for a project.
func (c *RPSStatsCache) SetPrevBatchRPS(projectID domain.ProjectID, rps uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prevBatchRPS[projectID] = rps
}

// GetPrevBatchRPS gets the previous batch RPS for a project.
func (c *RPSStatsCache) GetPrevBatchRPS(projectID domain.ProjectID) (uint64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	rps, ok := c.prevBatchRPS[projectID]

	return rps, ok
}

// SetCurrentBatch sets the current time batch for a project.
func (c *RPSStatsCache) SetCurrentBatch(projectID domain.ProjectID, batch int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.currentBatch[projectID] = batch
}

// GetCurrentBatch gets the current time batch for a project.
func (c *RPSStatsCache) GetCurrentBatch(projectID domain.ProjectID) (int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	batch, ok := c.currentBatch[projectID]

	return batch, ok
}

// CalculateRPS calculates the RPS based on the count and window size.
func CalculateRPS(count uint64, windowSize time.Duration) uint64 {
	return uint64(float64(count) / windowSize.Seconds())
}

// GetThrottleEnabled returns whether throttle is enabled for a project.
func (c *RPSStatsCache) GetThrottleEnabled(projectID domain.ProjectID) (bool, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	enabled, ok := c.throttleEnabled[projectID]

	return enabled, ok
}

// SetThrottleEnabled sets whether throttle is enabled for a project.
func (c *RPSStatsCache) SetThrottleEnabled(projectID domain.ProjectID, enabled bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.throttleEnabled[projectID] = enabled
}

// Get returns the RPS stats for a project (for backward compatibility).
func (c *RPSStatsCache) Get(projectID domain.ProjectID) (RPSStats, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	rps, rpsOk := c.rpsStats[projectID]
	enabled, enabledOk := c.throttleEnabled[projectID]

	if !rpsOk && !enabledOk {
		return RPSStats{}, false
	}

	return RPSStats{
		RPS:             rps,
		ThrottleEnabled: enabled,
	}, true
}

// Set sets the RPS stats for a project (for backward compatibility).
func (c *RPSStatsCache) Set(projectID domain.ProjectID, stats RPSStats) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rpsStats[projectID] = stats.RPS
	c.throttleEnabled[projectID] = stats.ThrottleEnabled
}

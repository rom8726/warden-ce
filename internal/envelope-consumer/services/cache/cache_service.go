//nolint:nestif // need refactoring
package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	commonconfig "github.com/rom8726/warden/internal/common/config"
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/envelope-consumer/contract"
	"github.com/rom8726/warden/pkg/metrics"
)

// Service provides caching functionality for envelope/store operations.
type Service struct {
	cacheManager contract.CacheManager
	config       commonconfig.CacheConfig
	mu           sync.RWMutex
	closed       bool
}

// Ensure Service implements contract.CacheService.
var _ contract.CacheService = (*Service)(nil)

// New creates a new cache service.
func New(config *commonconfig.CacheConfig, cacheManager contract.CacheManager) (*Service, error) {
	if !config.Enabled {
		slog.Info("Cache is disabled")

		return &Service{
			config: *config,
			closed: true,
		}, nil
	}

	if cacheManager == nil {
		return nil, errors.New("cache manager is required when cache is enabled")
	}

	// Initialize metrics
	metrics.CacheCapacity.WithLabelValues("release").Set(float64(config.ReleaseCacheSize))
	metrics.CacheCapacity.WithLabelValues("issue").Set(float64(config.IssueCacheSize))
	metrics.CacheCapacity.WithLabelValues("issue_release").Set(float64(config.IssueReleaseCacheSize))

	slog.Info("Cache service initialized",
		"release_cache_size", config.ReleaseCacheSize,
		"issue_cache_size", config.IssueCacheSize,
		"issue_release_cache_size", config.IssueReleaseCacheSize,
	)

	return &Service{
		cacheManager: cacheManager,
		config:       *config,
	}, nil
}

// GetOrCreateRelease retrieves a release from cache or creates it via repository.
func (s *Service) GetOrCreateRelease(
	ctx context.Context,
	projectID domain.ProjectID,
	version string,
	releaseRepo contract.ReleaseRepository,
) (domain.ReleaseID, error) {
	if s.closed {
		// Cache is disabled, create directly
		releaseID, err := releaseRepo.Create(ctx, domain.ReleaseDTO{
			ProjectID:   projectID,
			Version:     version,
			Description: "",
		})
		if err != nil {
			return 0, fmt.Errorf("create release: %w", err)
		}

		return releaseID, nil
	}

	// Try to get from the cache first
	if releaseValue, found := s.cacheManager.GetRelease(ctx, uint(projectID), version); found {
		metrics.CacheHits.WithLabelValues("release").Inc()
		slog.Debug("Release found in cache",
			"project_id", projectID,
			"version", version,
			"release_id", releaseValue.ReleaseID,
		)

		return domain.ReleaseID(releaseValue.ReleaseID), nil
	}

	// Cache miss
	metrics.CacheMisses.WithLabelValues("release").Inc()

	// Cache miss, create in a database
	releaseID, err := releaseRepo.Create(ctx, domain.ReleaseDTO{
		ProjectID:   projectID,
		Version:     version,
		Description: "",
	})
	if err != nil {
		return 0, fmt.Errorf("create release: %w", err)
	}

	// Store in cache
	if err := s.cacheManager.SetRelease(ctx, uint(projectID), version, uint(releaseID)); err != nil {
		slog.Warn("Failed to store release in cache",
			"error", err,
			"project_id", projectID,
			"version", version,
		)
	} else {
		slog.Debug("Release stored in cache",
			"project_id", projectID,
			"version", version,
			"release_id", releaseID,
		)
		// Update cache size metric
		stats := s.cacheManager.Stats()
		if releaseStats, ok := stats["release_cache"].(map[string]any); ok {
			if size, ok := releaseStats["size"].(int); ok {
				metrics.CacheSize.WithLabelValues("release").Set(float64(size))
			}
		}
	}

	return releaseID, nil
}

// GetOrCreateIssue retrieves an issue from the cache or creates it via repository.
func (s *Service) GetOrCreateIssue(
	ctx context.Context,
	issue domain.IssueDTO,
	issueRepo contract.IssuesRepository,
) (domain.IssueUpsertResult, error) {
	if s.closed {
		// Cache is disabled, create directly
		return issueRepo.UpsertIssue(ctx, issue)
	}

	// Try to get from the cache first
	if issueValue, found := s.cacheManager.GetIssue(ctx, issue.Fingerprint); found {
		metrics.CacheHits.WithLabelValues("issue").Inc()
		slog.Debug("Issue found in cache",
			"fingerprint", issue.Fingerprint,
			"issue_id", issueValue.IssueID,
		)
		// Note: We still need to call UpsertIssue to update last_seen and total_events,
		// but we can optimize this in the future by storing more data in cache
		return issueRepo.UpsertIssue(ctx, issue)
	}

	// Cache miss
	metrics.CacheMisses.WithLabelValues("issue").Inc()

	// Cache miss, create in a database
	upsertResult, err := issueRepo.UpsertIssue(ctx, issue)
	if err != nil {
		return domain.IssueUpsertResult{}, fmt.Errorf("upsert issue: %w", err)
	}

	// Store in cache
	if err := s.cacheManager.SetIssue(ctx, issue.Fingerprint, uint(upsertResult.ID)); err != nil {
		slog.Warn("Failed to store issue in cache",
			"error", err,
			"fingerprint", issue.Fingerprint,
		)
	} else {
		slog.Debug("Issue stored in cache",
			"fingerprint", issue.Fingerprint,
			"issue_id", upsertResult.ID,
		)
		// Update cache size metric
		stats := s.cacheManager.Stats()
		if issueStats, ok := stats["issue_cache"].(map[string]any); ok {
			if size, ok := issueStats["size"].(int); ok {
				metrics.CacheSize.WithLabelValues("issue").Set(float64(size))
			}
		}
	}

	return upsertResult, nil
}

// GetOrCreateIssueRelease retrieves an issue_release from the cache or creates it via repository.
func (s *Service) GetOrCreateIssueRelease(
	ctx context.Context,
	issueID domain.IssueID,
	releaseID domain.ReleaseID,
	firstSeenIn bool,
	issueReleaseRepo contract.IssueReleasesRepository,
) error {
	if s.closed {
		// Cache is disabled, create directly
		return issueReleaseRepo.Create(ctx, issueID, releaseID, firstSeenIn)
	}

	// Try to get from the cache first
	if _, found := s.cacheManager.GetIssueRelease(ctx, uint(issueID), uint(releaseID)); found {
		metrics.CacheHits.WithLabelValues("issue_release").Inc()
		slog.Debug("Issue release found in cache",
			"issue_id", issueID,
			"release_id", releaseID,
		)
		// Issue release already exists, no need to create
		return nil
	}

	// Cache miss
	metrics.CacheMisses.WithLabelValues("issue_release").Inc()

	// Cache miss, create in a database
	if err := issueReleaseRepo.Create(ctx, issueID, releaseID, firstSeenIn); err != nil {
		return fmt.Errorf("create issue release: %w", err)
	}

	// Store in cache (we don't have issue_release ID, so we store with a fake ID)
	// In a real implementation, we might want to modify the repository to return the ID
	if err := s.cacheManager.SetIssueRelease(ctx, uint(issueID), uint(releaseID), 0, firstSeenIn); err != nil {
		slog.Warn("Failed to store issue release in cache",
			"error", err,
			"issue_id", issueID,
			"release_id", releaseID,
		)
	} else {
		slog.Debug("Issue release stored in cache",
			"issue_id", issueID,
			"release_id", releaseID,
		)
		// Update cache size metric
		stats := s.cacheManager.Stats()
		if issueReleaseStats, ok := stats["issue_release_cache"].(map[string]any); ok {
			if size, ok := issueReleaseStats["size"].(int); ok {
				metrics.CacheSize.WithLabelValues("issue_release").Set(float64(size))
			}
		}
	}

	return nil
}

// Stats return cache statistics.
func (s *Service) Stats() map[string]any {
	if s.closed {
		return map[string]any{
			"status": "disabled",
		}
	}

	return s.cacheManager.Stats()
}

// Clear clears all caches.
func (s *Service) Clear(ctx context.Context) error {
	if s.closed {
		return nil
	}

	// Reset cache size metrics
	metrics.CacheSize.WithLabelValues("release").Set(0)
	metrics.CacheSize.WithLabelValues("issue").Set(0)
	metrics.CacheSize.WithLabelValues("issue_release").Set(0)

	return s.cacheManager.Clear(ctx)
}

// Close closes the cache service.
func (s *Service) Close(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	// Reset cache size metrics
	metrics.CacheSize.WithLabelValues("release").Set(0)
	metrics.CacheSize.WithLabelValues("issue").Set(0)
	metrics.CacheSize.WithLabelValues("issue_release").Set(0)

	if s.cacheManager != nil {
		return s.cacheManager.Close(ctx)
	}

	return nil
}

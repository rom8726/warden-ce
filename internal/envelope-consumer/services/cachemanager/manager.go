package cachemanager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	commonconfig "github.com/rom8726/warden/internal/common/config"
	"github.com/rom8726/warden/internal/envelope-consumer/contract"
	"github.com/rom8726/warden/pkg/cache"
)

// Manager manages all caches.
type Manager struct {
	releaseCache      cache.Cache[ReleaseKey, ReleaseValue]
	issueCache        cache.Cache[IssueKey, IssueValue]
	issueReleaseCache cache.Cache[IssueReleaseKey, IssueReleaseValue]

	mu     sync.RWMutex
	closed bool
}

// Ensure Manager implements contract.CacheManager.
var _ contract.CacheManager = (*Manager)(nil)

// New creates a new cache manager.
func New(config *commonconfig.CacheConfig) (*Manager, error) {
	releaseCache, err := cache.NewLRUCache[ReleaseKey, ReleaseValue](config.ReleaseCacheSize)
	if err != nil {
		return nil, fmt.Errorf("create release cache: %w", err)
	}

	issueCache, err := cache.NewLRUCache[IssueKey, IssueValue](config.IssueCacheSize)
	if err != nil {
		return nil, fmt.Errorf("create issue cache: %w", err)
	}

	issueReleaseCache, err := cache.NewLRUCache[IssueReleaseKey, IssueReleaseValue](config.IssueReleaseCacheSize)
	if err != nil {
		return nil, fmt.Errorf("create issue release cache: %w", err)
	}

	return &Manager{
		releaseCache:      releaseCache,
		issueCache:        issueCache,
		issueReleaseCache: issueReleaseCache,
	}, nil
}

// ReleaseCache returns the release cache.
func (m *Manager) ReleaseCache() cache.Cache[ReleaseKey, ReleaseValue] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil
	}

	return m.releaseCache
}

// IssueCache returns the issue cache.
func (m *Manager) IssueCache() cache.Cache[IssueKey, IssueValue] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil
	}

	return m.issueCache
}

// IssueReleaseCache returns the issue release cache.
func (m *Manager) IssueReleaseCache() cache.Cache[IssueReleaseKey, IssueReleaseValue] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil
	}

	return m.issueReleaseCache
}

// GetRelease retrieves a release from a cache.
func (m *Manager) GetRelease(ctx context.Context, projectID uint, version string) (contract.ReleaseValue, bool) {
	releaseCache := m.ReleaseCache()
	if releaseCache == nil {
		return contract.ReleaseValue{}, false
	}

	key := ReleaseKey{
		ProjectID: projectID,
		Version:   version,
	}

	value, found := releaseCache.Get(ctx, key)

	return contract.ReleaseValue{
		ReleaseID: value.ReleaseID,
	}, found
}

// SetRelease stores a release in cache.
func (m *Manager) SetRelease(ctx context.Context, projectID uint, version string, releaseID uint) error {
	releaseCache := m.ReleaseCache()
	if releaseCache == nil {
		return errors.New("releaseCache manager is closed")
	}

	key := ReleaseKey{
		ProjectID: projectID,
		Version:   version,
	}

	value := ReleaseValue{
		ReleaseID: releaseID,
	}

	return releaseCache.Set(ctx, key, value)
}

// GetIssue retrieves an issue from cache.
func (m *Manager) GetIssue(ctx context.Context, fingerprint string) (contract.IssueValue, bool) {
	issueCache := m.IssueCache()
	if issueCache == nil {
		return contract.IssueValue{}, false
	}

	key := IssueKey{
		Fingerprint: fingerprint,
	}

	value, found := issueCache.Get(ctx, key)

	return contract.IssueValue{
		IssueID: value.IssueID,
	}, found
}

// SetIssue stores an issue in cache.
func (m *Manager) SetIssue(ctx context.Context, fingerprint string, issueID uint) error {
	issueCache := m.IssueCache()
	if issueCache == nil {
		return errors.New("issueCache manager is closed")
	}

	key := IssueKey{
		Fingerprint: fingerprint,
	}

	value := IssueValue{
		IssueID: issueID,
	}

	return issueCache.Set(ctx, key, value)
}

// GetIssueRelease retrieves an issue release from the cache.
func (m *Manager) GetIssueRelease(ctx context.Context, issueID, releaseID uint) (contract.IssueReleaseValue, bool) {
	issueReleaseCache := m.IssueReleaseCache()
	if issueReleaseCache == nil {
		return contract.IssueReleaseValue{}, false
	}

	key := IssueReleaseKey{
		IssueID:   issueID,
		ReleaseID: releaseID,
	}

	value, found := issueReleaseCache.Get(ctx, key)

	return contract.IssueReleaseValue{
		IssueReleaseID: value.IssueReleaseID,
		FirstSeenIn:    value.FirstSeenIn,
	}, found
}

// SetIssueRelease stores an issue release in a cache.
func (m *Manager) SetIssueRelease(
	ctx context.Context,
	issueID, releaseID, issueReleaseID uint,
	firstSeenIn bool,
) error {
	issueReleaseCache := m.IssueReleaseCache()
	if issueReleaseCache == nil {
		return errors.New("issueReleaseCache manager is closed")
	}

	key := IssueReleaseKey{
		IssueID:   issueID,
		ReleaseID: releaseID,
	}

	value := IssueReleaseValue{
		IssueReleaseID: issueReleaseID,
		FirstSeenIn:    firstSeenIn,
	}

	return issueReleaseCache.Set(ctx, key, value)
}

// Stats returns cache statistics.
func (m *Manager) Stats() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return map[string]any{
			"status": "closed",
		}
	}

	return map[string]any{
		"status": "active",
		"release_cache": map[string]any{
			"size":     m.releaseCache.Size(),
			"capacity": m.releaseCache.Capacity(),
		},
		"issue_cache": map[string]any{
			"size":     m.issueCache.Size(),
			"capacity": m.issueCache.Capacity(),
		},
		"issue_release_cache": map[string]any{
			"size":     m.issueReleaseCache.Size(),
			"capacity": m.issueReleaseCache.Capacity(),
		},
	}
}

// Clear clears all caches.
func (m *Manager) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return errors.New("cache manager is closed")
	}

	if err := m.releaseCache.Clear(ctx); err != nil {
		return fmt.Errorf("clear release cache: %w", err)
	}

	if err := m.issueCache.Clear(ctx); err != nil {
		return fmt.Errorf("clear issue cache: %w", err)
	}

	if err := m.issueReleaseCache.Clear(ctx); err != nil {
		return fmt.Errorf("clear issue release cache: %w", err)
	}

	slog.Info("All caches cleared")

	return nil
}

// Close closes the cache manager and clears all caches.
func (m *Manager) Close(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}

	m.closed = true

	if err := m.releaseCache.Clear(ctx); err != nil {
		slog.Error("Failed to clear release cache on close", "error", err)
	}

	if err := m.issueCache.Clear(ctx); err != nil {
		slog.Error("Failed to clear issue cache on close", "error", err)
	}

	if err := m.issueReleaseCache.Clear(ctx); err != nil {
		slog.Error("Failed to clear issue release cache on close", "error", err)
	}

	slog.Info("Cache manager closed")

	return nil
}

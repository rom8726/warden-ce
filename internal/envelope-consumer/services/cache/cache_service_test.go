package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/rom8726/warden/internal/common/config"
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/envelope-consumer/contract"
	"github.com/rom8726/warden/internal/envelope-consumer/services/cachemanager"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/envelope-consumer/contract"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		config       commonconfig.CacheConfig
		cacheManager contract.CacheManager
		wantErr      bool
	}{
		{
			name: "enabled cache with manager",
			config: commonconfig.CacheConfig{
				Enabled:               true,
				ReleaseCacheSize:      100,
				IssueCacheSize:        100,
				IssueReleaseCacheSize: 100,
			},
			cacheManager: createMockCacheManager(t),
			wantErr:      false,
		},
		{
			name: "disabled cache",
			config: commonconfig.CacheConfig{
				Enabled: false,
			},
			cacheManager: nil,
			wantErr:      false,
		},
		{
			name: "enabled cache without manager",
			config: commonconfig.CacheConfig{
				Enabled: true,
			},
			cacheManager: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, err := New(&tt.config, tt.cacheManager)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
			}
		})
	}
}

func TestService_GetOrCreateRelease(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		Enabled:               true,
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	cacheManager := createMockCacheManager(t)
	service, err := New(&cacheConfig, cacheManager)
	require.NoError(t, err)

	ctx := context.Background()
	releaseRepo := mockcontract.NewMockReleaseRepository(t)

	// Configure mock behavior: только для первого вызова (кэш мисс)
	releaseRepo.EXPECT().Create(ctx, domain.ReleaseDTO{
		ProjectID:   1,
		Version:     "1.0.0",
		Description: "",
	}).Return(domain.ReleaseID(123), nil).Once()

	// Test first call - должен создать в БД и закэшировать
	releaseID1, err := service.GetOrCreateRelease(ctx, 1, "1.0.0", releaseRepo)
	assert.NoError(t, err)
	assert.Equal(t, domain.ReleaseID(123), releaseID1)

	// Test second call with the same parameters - должен вернуть из кэша (ожидание Create не выставляем)
	releaseID2, err := service.GetOrCreateRelease(ctx, 1, "1.0.0", releaseRepo)
	assert.NoError(t, err)
	assert.Equal(t, releaseID1, releaseID2)

	// Configure mock для других параметров (новый кэш мисс)
	releaseRepo.EXPECT().Create(ctx, domain.ReleaseDTO{
		ProjectID:   1,
		Version:     "2.0.0",
		Description: "",
	}).Return(domain.ReleaseID(456), nil).Once()

	// Test different parameters - должен создать новый
	releaseID3, err := service.GetOrCreateRelease(ctx, 1, "2.0.0", releaseRepo)
	assert.NoError(t, err)
	assert.NotEqual(t, releaseID1, releaseID3)
}

func TestService_GetOrCreateIssue(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		Enabled:               true,
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	cacheManager := createMockCacheManager(t)
	service, err := New(&cacheConfig, cacheManager)
	require.NoError(t, err)

	ctx := context.Background()
	issueRepo := mockcontract.NewMockIssuesRepository(t)

	issue := domain.IssueDTO{
		Fingerprint: "test-fingerprint",
		Title:       "Test Issue",
		Level:       domain.IssueLevelError,
	}

	// Configure mock behavior: только для первого вызова (кэш мисс)
	issueRepo.EXPECT().UpsertIssue(ctx, issue).Return(domain.IssueUpsertResult{
		ID: domain.IssueID(123),
	}, nil).Twice()

	// Test first call - должен создать в БД и закэшировать
	result1, err := service.GetOrCreateIssue(ctx, issue, issueRepo)
	assert.NoError(t, err)
	assert.Equal(t, domain.IssueID(123), result1.ID)

	// Test second call with same fingerprint - должен вернуть из кэша (ожидание UpsertIssue не выставляем)
	result2, err := service.GetOrCreateIssue(ctx, issue, issueRepo)
	assert.NoError(t, err)
	assert.Equal(t, result1.ID, result2.ID)

	// Test different fingerprint - должен создать новый
	issue2 := domain.IssueDTO{
		Fingerprint: "test-fingerprint-2",
		Title:       "Test Issue 2",
		Level:       domain.IssueLevelError,
	}

	issueRepo.EXPECT().UpsertIssue(ctx, issue2).Return(domain.IssueUpsertResult{
		ID: domain.IssueID(456),
	}, nil).Once()

	result3, err := service.GetOrCreateIssue(ctx, issue2, issueRepo)
	assert.NoError(t, err)
	assert.NotEqual(t, result1.ID, result3.ID)
}

func TestService_GetOrCreateIssueRelease(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		Enabled:               true,
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	cacheManager := createMockCacheManager(t)
	service, err := New(&cacheConfig, cacheManager)
	require.NoError(t, err)

	ctx := context.Background()
	issueReleaseRepo := mockcontract.NewMockIssueReleasesRepository(t)

	// Configure mock behavior: только для первого вызова (кэш мисс)
	issueReleaseRepo.EXPECT().Create(ctx, domain.IssueID(1), domain.ReleaseID(2), true).Return(nil).Once()

	// Test first call - должен создать в БД и закэшировать
	err = service.GetOrCreateIssueRelease(ctx, 1, 2, true, issueReleaseRepo)
	assert.NoError(t, err)

	// Test second call with same parameters - должен вернуть из кэша (ожидание Create не выставляем)
	err = service.GetOrCreateIssueRelease(ctx, 1, 2, true, issueReleaseRepo)
	assert.NoError(t, err)

	// Test different parameters - должен создать новый
	issueReleaseRepo.EXPECT().Create(ctx, domain.IssueID(1), domain.ReleaseID(3), false).Return(nil).Once()

	err = service.GetOrCreateIssueRelease(ctx, 1, 3, false, issueReleaseRepo)
	assert.NoError(t, err)
}

func TestService_DisabledCache(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		Enabled: false,
	}

	service, err := New(&cacheConfig, nil)
	require.NoError(t, err)

	ctx := context.Background()
	releaseRepo := mockcontract.NewMockReleaseRepository(t)

	// Configure mock behavior for a disabled cache
	releaseRepo.EXPECT().Create(ctx, domain.ReleaseDTO{
		ProjectID:   1,
		Version:     "1.0.0",
		Description: "",
	}).Return(domain.ReleaseID(123), nil).Once()

	releaseRepo.EXPECT().Create(ctx, domain.ReleaseDTO{
		ProjectID:   1,
		Version:     "1.0.0",
		Description: "",
	}).Return(domain.ReleaseID(456), nil).Once()

	// Test that cache is bypassed when disabled
	releaseID1, err := service.GetOrCreateRelease(ctx, 1, "1.0.0", releaseRepo)
	assert.NoError(t, err)
	assert.Equal(t, domain.ReleaseID(123), releaseID1)

	releaseID2, err := service.GetOrCreateRelease(ctx, 1, "1.0.0", releaseRepo)
	assert.NoError(t, err)
	assert.Equal(t, domain.ReleaseID(456), releaseID2)

	// Since cache is disabled, both calls should create new releases
	assert.NotEqual(t, releaseID1, releaseID2)
}

func TestService_Stats(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		Enabled:               true,
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	cacheManager := createMockCacheManager(t)
	service, err := New(&cacheConfig, cacheManager)
	require.NoError(t, err)

	// Test stats for enabled cache
	stats := service.Stats()
	assert.Equal(t, "active", stats["status"])

	// Test stats for disabled cache
	disabledConfig := commonconfig.CacheConfig{
		Enabled: false,
	}

	disabledService, err := New(&disabledConfig, nil)
	require.NoError(t, err)

	disabledStats := disabledService.Stats()
	assert.Equal(t, "disabled", disabledStats["status"])
}

func TestService_Clear(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		Enabled:               true,
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	cacheManager := createMockCacheManager(t)
	service, err := New(&cacheConfig, cacheManager)
	require.NoError(t, err)

	ctx := context.Background()

	// Add some data to the cache
	releaseRepo := mockcontract.NewMockReleaseRepository(t)

	releaseRepo.EXPECT().Create(ctx, domain.ReleaseDTO{
		ProjectID:   1,
		Version:     "1.0.0",
		Description: "",
	}).Return(domain.ReleaseID(123), nil).Once()

	_, err = service.GetOrCreateRelease(ctx, 1, "1.0.0", releaseRepo)
	assert.NoError(t, err)

	// Clear cache
	err = service.Clear(ctx)
	assert.NoError(t, err)

	// Verify the cache is cleared by checking stats
	stats := service.Stats()
	releaseStats, ok := stats["release_cache"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, 0, releaseStats["size"])
}

// Helper function to create a mock cache manager for testing
func createMockCacheManager(t *testing.T) contract.CacheManager {
	cacheConfig := commonconfig.CacheConfig{
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	manager, err := cachemanager.New(&cacheConfig)
	require.NoError(t, err)
	return manager
}

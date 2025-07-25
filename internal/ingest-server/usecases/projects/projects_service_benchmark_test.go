package projects

import (
	"context"
	"testing"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/ingest-server/contract"
)

// mockProjectsRepository is a simple mock for benchmarking
type mockProjectsRepository struct {
	callCount int
}

func (m *mockProjectsRepository) GetProjectIDs(ctx context.Context) ([]domain.ProjectID, error) {
	return nil, nil
}

func (m *mockProjectsRepository) ValidateProjectKey(ctx context.Context, projectID domain.ProjectID, key string) (bool, error) {
	m.callCount++
	// Simulate successful validation for benchmarking
	return true, nil
}

// ProjectServiceWithoutCache is a version without caching for benchmarking
type ProjectServiceWithoutCache struct {
	projectRepo contract.ProjectsRepository
}

func (s *ProjectServiceWithoutCache) ValidateProjectKey(
	ctx context.Context,
	projectID domain.ProjectID,
	key string,
) (bool, error) {
	return s.projectRepo.ValidateProjectKey(ctx, projectID, key)
}

func BenchmarkValidateProjectKey_WithCache(b *testing.B) {
	// Create a mock repository
	mockRepo := &mockProjectsRepository{}

	// Create service with cache
	service := &ProjectService{
		projectRepo: mockRepo,
		keyCache:    make(map[ProjectCacheKey]struct{}),
	}

	ctx := context.Background()
	projectID := domain.ProjectID(1)
	key := "valid-key"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateProjectKey(ctx, projectID, key)
	}

	// Print cache hit ratio
	b.ReportMetric(float64(mockRepo.callCount), "repo_calls")
}

func BenchmarkValidateProjectKey_WithoutCache(b *testing.B) {
	// Create a mock repository
	mockRepo := &mockProjectsRepository{}

	// Create service without cache
	service := &ProjectServiceWithoutCache{
		projectRepo: mockRepo,
	}

	ctx := context.Background()
	projectID := domain.ProjectID(1)
	key := "valid-key"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateProjectKey(ctx, projectID, key)
	}

	// Print cache hit ratio
	b.ReportMetric(float64(mockRepo.callCount), "repo_calls")
}

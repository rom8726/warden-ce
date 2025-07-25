package projects

import (
	"context"
	"sync"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/ingest-server/contract"
)

type ProjectService struct {
	projectRepo contract.ProjectsRepository
	cacheMu     sync.RWMutex
	keyCache    map[ProjectCacheKey]struct{}
}

func New(
	projectRepo contract.ProjectsRepository,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		keyCache:    make(map[ProjectCacheKey]struct{}),
	}
}

func (s *ProjectService) ValidateProjectKey(
	ctx context.Context,
	projectID domain.ProjectID,
	key string,
) (bool, error) {
	cacheKey := ProjectCacheKey{
		ProjectID: projectID,
		Key:       key,
	}

	s.cacheMu.RLock()
	_, found := s.keyCache[cacheKey]
	s.cacheMu.RUnlock()
	if found {
		return true, nil
	}

	allowed, err := s.projectRepo.ValidateProjectKey(ctx, projectID, key)
	if err != nil {
		return false, err
	}

	if allowed {
		s.cacheMu.Lock()
		s.keyCache[cacheKey] = struct{}{}
		s.cacheMu.Unlock()
	}

	return allowed, nil
}

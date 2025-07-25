package projects

import (
	"fmt"

	"github.com/rom8726/warden/internal/domain"
)

type ProjectCacheKey struct {
	ProjectID domain.ProjectID
	Key       string
}

func (k ProjectCacheKey) String() string {
	return fmt.Sprintf("%d:%s", k.ProjectID, k.Key)
}

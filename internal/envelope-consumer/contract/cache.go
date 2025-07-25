package contract

import (
	"context"

	"github.com/rom8726/warden/internal/domain"
)

// CacheManager defines the interface for cache management operations.
type CacheManager interface {
	GetRelease(ctx context.Context, projectID uint, version string) (ReleaseValue, bool)
	SetRelease(ctx context.Context, projectID uint, version string, releaseID uint) error

	GetIssue(ctx context.Context, fingerprint string) (IssueValue, bool)
	SetIssue(ctx context.Context, fingerprint string, issueID uint) error

	GetIssueRelease(ctx context.Context, issueID, releaseID uint) (IssueReleaseValue, bool)
	SetIssueRelease(ctx context.Context, issueID, releaseID, issueReleaseID uint, firstSeenIn bool) error

	Stats() map[string]any
	Clear(ctx context.Context) error
	Close(ctx context.Context) error
}

// CacheService defines the interface for cache service operations.
type CacheService interface {
	// GetOrCreateRelease retrieves a release from cache or creates it via repository
	GetOrCreateRelease(
		ctx context.Context,
		projectID domain.ProjectID,
		version string,
		releaseRepo ReleaseRepository,
	) (domain.ReleaseID, error)

	// GetOrCreateIssue retrieves an issue from cache or creates it via repository
	GetOrCreateIssue(
		ctx context.Context,
		issue domain.IssueDTO,
		issueRepo IssuesRepository,
	) (domain.IssueUpsertResult, error)

	// GetOrCreateIssueRelease retrieves an issue_release from cache or creates it via repository
	GetOrCreateIssueRelease(
		ctx context.Context,
		issueID domain.IssueID,
		releaseID domain.ReleaseID,
		firstSeenIn bool,
		issueReleaseRepo IssueReleasesRepository,
	) error

	Stats() map[string]any
	Clear(ctx context.Context) error
	Close(ctx context.Context) error
}

type ReleaseValue struct {
	ReleaseID uint
}

type IssueValue struct {
	IssueID uint
}

type IssueReleaseValue struct {
	IssueReleaseID uint
	FirstSeenIn    bool
}

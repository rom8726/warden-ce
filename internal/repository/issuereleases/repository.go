package issuereleases

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/db"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

func (r *Repository) Create(
	ctx context.Context,
	issueID domain.IssueID,
	releaseID domain.ReleaseID,
	firstSeenIn bool,
) error {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO issue_releases (issue_id, release_id, first_seen_in)
VALUES ($1, $2, $3)
ON CONFLICT (issue_id, release_id) DO NOTHING
`
	_, err := executor.Exec(ctx, query, issueID, releaseID, firstSeenIn)
	if err != nil {
		return fmt.Errorf("insert issue_release: %w", err)
	}

	return nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}

package teams

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db: pool,
	}
}

func (r *Repository) Create(ctx context.Context, teamDTO domain.TeamDTO) (domain.Team, error) {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO teams (name, created_at)
VALUES ($1, $2)
RETURNING *`

	rows, err := executor.Query(ctx, query,
		teamDTO.Name,
		time.Now(),
	)
	if err != nil {
		return domain.Team{}, fmt.Errorf("insert team: %w", err)
	}
	defer rows.Close()

	team, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[teamModel])
	if err != nil {
		return domain.Team{}, fmt.Errorf("collect team: %w", err)
	}

	return team.toDomain(), nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.TeamID) (domain.Team, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM teams WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Team{}, fmt.Errorf("query team by ID: %w", err)
	}
	defer rows.Close()

	team, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[teamModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, domain.ErrEntityNotFound
		}

		return domain.Team{}, fmt.Errorf("collect team: %w", err)
	}

	// Get team members
	members, err := r.GetMembers(ctx, id)
	if err != nil {
		return domain.Team{}, fmt.Errorf("get team members: %w", err)
	}

	result := team.toDomain()
	result.Members = members

	return result, nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (domain.Team, error) {
	executor := r.getExecutor(ctx)
	const query = `SELECT * FROM teams WHERE name = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, name)
	if err != nil {
		return domain.Team{}, fmt.Errorf("query team by name: %w", err)
	}
	defer rows.Close()

	team, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[teamModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, domain.ErrEntityNotFound
		}

		return domain.Team{}, fmt.Errorf("collect team: %w", err)
	}

	// Get team members
	members, err := r.GetMembers(ctx, domain.TeamID(team.ID))
	if err != nil {
		return domain.Team{}, fmt.Errorf("get team members: %w", err)
	}

	result := team.toDomain()
	result.Members = members

	return team.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, team domain.Team) error {
	executor := r.getExecutor(ctx)

	const query = `
UPDATE teams
SET name = $1
WHERE id = $2`

	_, err := executor.Exec(ctx, query,
		team.Name,
		team.ID,
	)
	if err != nil {
		return fmt.Errorf("update team: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id domain.TeamID) error {
	executor := r.getExecutor(ctx)

	// First delete all team members
	const deleteMembers = `
DELETE FROM team_members
WHERE team_id = $1`

	_, err := executor.Exec(ctx, deleteMembers, id)
	if err != nil {
		return fmt.Errorf("delete team members: %w", err)
	}

	// Then delete the team
	const deleteTeam = `
DELETE FROM teams
WHERE id = $1`

	_, err = executor.Exec(ctx, deleteTeam, id)
	if err != nil {
		return fmt.Errorf("delete team: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Team, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM teams ORDER BY id`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query teams: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[teamModel])
	if err != nil {
		return nil, fmt.Errorf("collect teams: %w", err)
	}

	teams := make([]domain.Team, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		teams = append(teams, model.toDomain())
	}

	// Get members for each team
	for i, team := range teams {
		members, err := r.GetMembers(ctx, team.ID)
		if err != nil {
			return nil, fmt.Errorf("get members for team %d: %w", team.ID, err)
		}
		teams[i].Members = members
	}

	return teams, nil
}

func (r *Repository) AddMember(
	ctx context.Context,
	teamID domain.TeamID,
	userID domain.UserID,
	role domain.Role,
) error {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO team_members (team_id, user_id, role)
VALUES ($1, $2, $3)`

	_, err := executor.Exec(ctx, query,
		teamID,
		userID,
		role,
	)
	if err != nil {
		return fmt.Errorf("insert team member: %w", err)
	}

	return nil
}

func (r *Repository) RemoveMember(ctx context.Context, teamID domain.TeamID, userID domain.UserID) error {
	executor := r.getExecutor(ctx)

	const query = `
DELETE FROM team_members
WHERE team_id = $1 AND user_id = $2`

	_, err := executor.Exec(ctx, query, teamID, userID)
	if err != nil {
		return fmt.Errorf("delete team member: %w", err)
	}

	return nil
}

func (r *Repository) UpdateMemberRole(
	ctx context.Context,
	teamID domain.TeamID,
	userID domain.UserID,
	newRole domain.Role,
) error {
	executor := r.getExecutor(ctx)

	const query = `
UPDATE team_members
SET role = $1
WHERE team_id = $2 AND user_id = $3`

	result, err := executor.Exec(ctx, query, string(newRole), teamID, userID)
	if err != nil {
		return fmt.Errorf("update member role: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrEntityNotFound
	}

	return nil
}

func (r *Repository) GetMembers(ctx context.Context, teamID domain.TeamID) ([]domain.TeamMember, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM team_members WHERE team_id = $1`

	rows, err := executor.Query(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("query team members: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[teamMemberModel])
	if err != nil {
		return nil, fmt.Errorf("collect team members: %w", err)
	}

	members := make([]domain.TeamMember, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		members = append(members, model.toDomain())
	}

	return members, nil
}

func (r *Repository) GetUniqueUserIDsByTeamIDs(
	ctx context.Context,
	teamIDs []domain.TeamID,
) ([]domain.UserID, error) {
	executor := r.getExecutor(ctx)
	const query = `
SELECT DISTINCT tm.user_id
FROM team_members tm
WHERE tm.team_id = ANY($1)`
	rows, err := executor.Query(ctx, query, teamIDs)
	if err != nil {
		return nil, fmt.Errorf("query unique user IDs by team IDs: %w", err)
	}
	defer rows.Close()

	var ids []domain.UserID
	for rows.Next() {
		var id domain.UserID
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return ids, nil
}

func (r *Repository) GetTeamsByUserID(ctx context.Context, userID domain.UserID) ([]domain.Team, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT t.*
FROM teams t
JOIN team_members tm ON t.id = tm.team_id
WHERE tm.user_id = $1`

	rows, err := executor.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query teams by user ID: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[teamModel])
	if err != nil {
		return nil, fmt.Errorf("collect teams: %w", err)
	}

	teams := make([]domain.Team, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		teams = append(teams, model.toDomain())
	}

	// Get members for each team
	for i, team := range teams {
		members, err := r.GetMembers(ctx, team.ID)
		if err != nil {
			return nil, fmt.Errorf("get members for team %d: %w", team.ID, err)
		}
		teams[i].Members = members
	}

	return teams, nil
}

func (r *Repository) GetTeamsByUserIDs(
	ctx context.Context,
	userIDs []domain.UserID,
) (map[domain.UserID][]domain.Team, error) {
	executor := r.getExecutor(ctx)
	const query = `
SELECT t.*
FROM teams t
JOIN team_members tm ON t.id = tm.team_id
WHERE tm.user_id = ANY($1)`

	rows, err := executor.Query(ctx, query, userIDs)
	if err != nil {
		return nil, fmt.Errorf("query teams by user IDs: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[teamModel])
	if err != nil {
		return nil, fmt.Errorf("collect teams: %w", err)
	}

	userTeamIDs := make(map[domain.UserID][]domain.TeamID)

	teams := make([]domain.Team, 0, len(listModels))
	for _, model := range listModels {
		team := model.toDomain()
		members, err := r.GetMembers(ctx, team.ID)
		if err != nil {
			return nil, fmt.Errorf("get members for team %d: %w", team.ID, err)
		}

		team.Members = members
		teams = append(teams, team)

		for _, member := range members {
			userTeamIDs[member.UserID] = append(userTeamIDs[member.UserID], team.ID)
		}
	}

	teamsByUserID := make(map[domain.UserID][]domain.Team)
	for _, userID := range userIDs {
		teamsIDs := userTeamIDs[userID]
		list := make([]domain.Team, 0, len(teamsIDs))
		for _, teamID := range teamsIDs {
			for _, team := range teams {
				if team.ID == teamID {
					list = append(list, team)
				}
			}
		}

		teamsByUserID[userID] = list
	}

	return teamsByUserID, nil
}

func (r *Repository) GetByProjectID(ctx context.Context, projectID domain.ProjectID) (domain.Team, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT t.*
FROM teams t
JOIN projects p ON p.team_id = t.id
WHERE p.id = $1
LIMIT 1`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return domain.Team{}, fmt.Errorf("query team by project ID: %w", err)
	}
	defer rows.Close()

	team, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[teamModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, domain.ErrEntityNotFound
		}

		return domain.Team{}, fmt.Errorf("collect team: %w", err)
	}

	// Get team members
	members, err := r.GetMembers(ctx, domain.TeamID(team.ID))
	if err != nil {
		return domain.Team{}, fmt.Errorf("get team members: %w", err)
	}

	result := team.toDomain()
	result.Members = members

	return result, nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}

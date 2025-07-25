package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/go-faster/errors"
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

func (r *Repository) CreateSetting(
	ctx context.Context,
	settingDTO domain.NotificationSettingDTO,
) (domain.NotificationSetting, error) {
	executor := r.getExecutor(ctx)

	model := settingFromDTO(settingDTO)

	const query = `
INSERT INTO notification_settings (project_id, type, config, enabled, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *`

	rows, err := executor.Query(ctx, query,
		model.ProjectID,
		model.Type,
		model.Config,
		model.Enabled,
		model.CreatedAt,
		model.UpdatedAt,
	)
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("insert notification setting: %w", err)
	}
	defer rows.Close()

	setting, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[notificationSettingModel])
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("collect notification setting: %w", err)
	}

	return setting.toDomain(), nil
}

func (r *Repository) GetSettingByID(
	ctx context.Context,
	id domain.NotificationSettingID,
) (domain.NotificationSetting, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM notification_settings WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("query notification setting by ID: %w", err)
	}
	defer rows.Close()

	setting, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[notificationSettingModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.NotificationSetting{}, domain.ErrEntityNotFound
		}

		return domain.NotificationSetting{}, fmt.Errorf("collect notification setting: %w", err)
	}

	// Get rules for this setting
	rules, err := r.ListRules(ctx, domain.NotificationSettingID(setting.ID))
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("get notification rules: %w", err)
	}

	result := setting.toDomain()
	result.Rules = rules

	return result, nil
}

func (r *Repository) UpdateSetting(ctx context.Context, setting domain.NotificationSetting) error {
	executor := r.getExecutor(ctx)

	model := settingFromDomain(setting)
	model.UpdatedAt = time.Now() // Always update the updated_at timestamp

	const query = `
UPDATE notification_settings
SET project_id = $1, type = $2, config = $3, enabled = $4, updated_at = $5
WHERE id = $6`

	_, err := executor.Exec(ctx, query,
		model.ProjectID,
		model.Type,
		model.Config,
		model.Enabled,
		model.UpdatedAt,
		model.ID,
	)
	if err != nil {
		return fmt.Errorf("update notification setting: %w", err)
	}

	return nil
}

func (r *Repository) DeleteSetting(ctx context.Context, id domain.NotificationSettingID) error {
	executor := r.getExecutor(ctx)

	// First delete all rules for this setting
	const deleteRules = `
DELETE FROM notification_rules
WHERE notification_setting_id = $1`

	_, err := executor.Exec(ctx, deleteRules, id)
	if err != nil {
		return fmt.Errorf("delete notification rules: %w", err)
	}

	// Then delete the setting
	const deleteSetting = `
DELETE FROM notification_settings
WHERE id = $1`

	_, err = executor.Exec(ctx, deleteSetting, id)
	if err != nil {
		return fmt.Errorf("delete notification setting: %w", err)
	}

	return nil
}

func (r *Repository) ListSettings(
	ctx context.Context,
	projectID domain.ProjectID,
) ([]domain.NotificationSetting, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT *
FROM notification_settings
WHERE project_id = $1
ORDER BY id`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("query notification settings: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[notificationSettingModel])
	if err != nil {
		return nil, fmt.Errorf("collect notification settings: %w", err)
	}

	settings := make([]domain.NotificationSetting, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		settings = append(settings, model.toDomain())
	}

	// Get rules for each setting
	for i, setting := range settings {
		rules, err := r.ListRules(ctx, setting.ID)
		if err != nil {
			return nil, fmt.Errorf("get rules for setting %d: %w", setting.ID, err)
		}
		settings[i].Rules = rules
	}

	return settings, nil
}

func (r *Repository) ListAllSettings(
	ctx context.Context,
) ([]domain.NotificationSetting, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT *
FROM notification_settings
ORDER BY id`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query notification settings: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[notificationSettingModel])
	if err != nil {
		return nil, fmt.Errorf("collect notification settings: %w", err)
	}

	settings := make([]domain.NotificationSetting, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		settings = append(settings, model.toDomain())
	}

	// Get rules for each setting
	for i, setting := range settings {
		rules, err := r.ListRules(ctx, setting.ID)
		if err != nil {
			return nil, fmt.Errorf("get rules for setting %d: %w", setting.ID, err)
		}
		settings[i].Rules = rules
	}

	return settings, nil
}

func (r *Repository) CreateRule(
	ctx context.Context,
	ruleDTO domain.NotificationRuleDTO,
) (domain.NotificationRule, error) {
	executor := r.getExecutor(ctx)

	model := ruleFromDTO(ruleDTO)

	const query = `
INSERT INTO notification_rules 
(notification_setting_id, event_level, fingerprint, is_new_error, is_regression, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *`

	rows, err := executor.Query(ctx, query,
		model.NotificationSetting,
		model.EventLevel,
		model.Fingerprint,
		model.IsNewError,
		model.IsRegression,
		model.CreatedAt,
	)
	if err != nil {
		return domain.NotificationRule{}, fmt.Errorf("insert notification rule: %w", err)
	}
	defer rows.Close()

	rule, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[notificationRuleModel])
	if err != nil {
		return domain.NotificationRule{}, fmt.Errorf("collect notification rule: %w", err)
	}

	return rule.toDomain(), nil
}

func (r *Repository) GetRuleByID(ctx context.Context, id domain.NotificationRuleID) (domain.NotificationRule, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM notification_rules WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.NotificationRule{}, fmt.Errorf("query notification rule by ID: %w", err)
	}
	defer rows.Close()

	rule, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[notificationRuleModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.NotificationRule{}, domain.ErrEntityNotFound
		}

		return domain.NotificationRule{}, fmt.Errorf("collect notification rule: %w", err)
	}

	return rule.toDomain(), nil
}

func (r *Repository) UpdateRule(ctx context.Context, rule domain.NotificationRule) error {
	executor := r.getExecutor(ctx)

	model := ruleFromDomain(rule)

	const query = `
UPDATE notification_rules
SET notification_setting_id = $1, event_level = $2, fingerprint = $3, is_new_error = $4, is_regression = $5
WHERE id = $6`

	_, err := executor.Exec(ctx, query,
		model.NotificationSetting,
		model.EventLevel,
		model.Fingerprint,
		model.IsNewError,
		model.IsRegression,
		model.ID,
	)
	if err != nil {
		return fmt.Errorf("update notification rule: %w", err)
	}

	return nil
}

func (r *Repository) DeleteRule(ctx context.Context, id domain.NotificationRuleID) error {
	executor := r.getExecutor(ctx)

	const query = `
DELETE FROM notification_rules
WHERE id = $1`

	_, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete notification rule: %w", err)
	}

	return nil
}

func (r *Repository) ListRules(
	ctx context.Context,
	settingID domain.NotificationSettingID,
) ([]domain.NotificationRule, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT *
FROM notification_rules
WHERE notification_setting_id = $1
ORDER BY id`

	rows, err := executor.Query(ctx, query, settingID)
	if err != nil {
		return nil, fmt.Errorf("query notification rules: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[notificationRuleModel])
	if err != nil {
		return nil, fmt.Errorf("collect notification rules: %w", err)
	}

	rules := make([]domain.NotificationRule, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		rules = append(rules, model.toDomain())
	}

	return rules, nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}

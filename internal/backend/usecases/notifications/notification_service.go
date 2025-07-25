package notifications

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rom8726/warden/internal/backend/contract"
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/db"
)

type Service struct {
	txManager                db.TxManager
	notificationSettingsRepo contract.NotificationSettingsRepository
	notificationRulesRepo    contract.NotificationRulesRepository
	notificationsQueueRepo   contract.NotificationsQueueRepository
	projectsRepo             contract.ProjectsRepository
	issuesRepo               contract.IssuesRepository

	notificationChannels []contract.NotificationChannel
}

func New(
	txManager db.TxManager,
	notificationSettingsRepo contract.NotificationSettingsRepository,
	notificationRulesRepo contract.NotificationRulesRepository,
	notificationsQueueRepo contract.NotificationsQueueRepository,
	projectsRepo contract.ProjectsRepository,
	issuesRepo contract.IssuesRepository,
	notificationChannels []contract.NotificationChannel,
) *Service {
	return &Service{
		txManager:                txManager,
		notificationSettingsRepo: notificationSettingsRepo,
		notificationRulesRepo:    notificationRulesRepo,
		notificationsQueueRepo:   notificationsQueueRepo,
		projectsRepo:             projectsRepo,
		issuesRepo:               issuesRepo,
		notificationChannels:     notificationChannels,
	}
}

// CreateNotificationSetting creates a new notification setting.
func (s *Service) CreateNotificationSetting(
	ctx context.Context,
	settingDTO domain.NotificationSettingDTO,
) (domain.NotificationSetting, error) {
	if _, err := s.projectsRepo.GetByID(ctx, settingDTO.ProjectID); err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("get project by ID: %w", err)
	}

	if settingDTO.Type == domain.NotificationTypeEmail {
		list, err := s.notificationSettingsRepo.ListSettings(ctx, settingDTO.ProjectID)
		if err != nil {
			return domain.NotificationSetting{}, fmt.Errorf("list notification settings: %w", err)
		}

		for _, setting := range list {
			if setting.Type == domain.NotificationTypeEmail {
				return domain.NotificationSetting{}, errors.New("email notification already exists")
			}
		}
	}

	result, err := s.notificationSettingsRepo.CreateSetting(ctx, settingDTO)
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("create notification setting: %w", err)
	}

	return result, nil
}

// GetNotificationSetting gets a notification setting by ID.
func (s *Service) GetNotificationSetting(
	ctx context.Context,
	id domain.NotificationSettingID,
) (domain.NotificationSetting, error) {
	setting, err := s.notificationSettingsRepo.GetSettingByID(ctx, id)
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("get notification setting: %w", err)
	}

	return setting, nil
}

// UpdateNotificationSetting updates a notification setting.
func (s *Service) UpdateNotificationSetting(
	ctx context.Context,
	setting domain.NotificationSetting,
) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if _, err := s.notificationSettingsRepo.GetSettingByID(ctx, setting.ID); err != nil {
			return fmt.Errorf("get notification setting: %w", err)
		}

		err := s.notificationSettingsRepo.UpdateSetting(ctx, setting)
		if err != nil {
			return fmt.Errorf("update notification setting: %w", err)
		}

		return nil
	})

	return err
}

// DeleteNotificationSetting deletes a notification setting.
func (s *Service) DeleteNotificationSetting(
	ctx context.Context,
	id domain.NotificationSettingID,
) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if _, err := s.notificationSettingsRepo.GetSettingByID(ctx, id); err != nil {
			return fmt.Errorf("get notification setting: %w", err)
		}

		err := s.notificationSettingsRepo.DeleteSetting(ctx, id)
		if err != nil {
			return fmt.Errorf("delete notification setting: %w", err)
		}

		return nil
	})

	return err
}

// ListNotificationSettings lists all notification settings for a project.
func (s *Service) ListNotificationSettings(
	ctx context.Context,
	projectID domain.ProjectID,
) ([]domain.NotificationSetting, error) {
	if _, err := s.projectsRepo.GetByID(ctx, projectID); err != nil {
		return nil, fmt.Errorf("get project by ID: %w", err)
	}

	settings, err := s.notificationSettingsRepo.ListSettings(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list notification settings: %w", err)
	}

	return settings, nil
}

// CreateNotificationRule creates a new notification rule.
func (s *Service) CreateNotificationRule(
	ctx context.Context,
	ruleDTO domain.NotificationRuleDTO,
) (domain.NotificationRule, error) {
	var result domain.NotificationRule
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if _, err := s.notificationSettingsRepo.GetSettingByID(ctx, ruleDTO.NotificationSetting); err != nil {
			return fmt.Errorf("get notification setting: %w", err)
		}

		var err error
		result, err = s.notificationRulesRepo.CreateRule(ctx, ruleDTO)
		if err != nil {
			return fmt.Errorf("create notification rule: %w", err)
		}

		return nil
	})
	if err != nil {
		return domain.NotificationRule{}, err
	}

	return result, nil
}

// GetNotificationRule gets a notification rule by ID.
func (s *Service) GetNotificationRule(
	ctx context.Context,
	id domain.NotificationRuleID,
) (domain.NotificationRule, error) {
	rule, err := s.notificationRulesRepo.GetRuleByID(ctx, id)
	if err != nil {
		return domain.NotificationRule{}, fmt.Errorf("get notification rule: %w", err)
	}

	return rule, nil
}

// UpdateNotificationRule updates a notification rule.
func (s *Service) UpdateNotificationRule(
	ctx context.Context,
	rule domain.NotificationRule,
) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if _, err := s.notificationSettingsRepo.GetSettingByID(ctx, rule.NotificationSetting); err != nil {
			return fmt.Errorf("get notification setting: %w", err)
		}

		err := s.notificationRulesRepo.UpdateRule(ctx, rule)
		if err != nil {
			return fmt.Errorf("update notification rule: %w", err)
		}

		return nil
	})

	return err
}

// DeleteNotificationRule deletes a notification rule.
func (s *Service) DeleteNotificationRule(
	ctx context.Context,
	id domain.NotificationRuleID,
) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if _, err := s.notificationRulesRepo.GetRuleByID(ctx, id); err != nil {
			return fmt.Errorf("get notification rule: %w", err)
		}

		err := s.notificationRulesRepo.DeleteRule(ctx, id)
		if err != nil {
			return fmt.Errorf("delete notification rule: %w", err)
		}

		return nil
	})

	return err
}

// ListNotificationRules lists all notification rules for a notification setting.
func (s *Service) ListNotificationRules(
	ctx context.Context,
	settingID domain.NotificationSettingID,
) ([]domain.NotificationRule, error) {
	if _, err := s.notificationSettingsRepo.GetSettingByID(ctx, settingID); err != nil {
		return nil, fmt.Errorf("get notification setting: %w", err)
	}

	rules, err := s.notificationRulesRepo.ListRules(ctx, settingID)
	if err != nil {
		return nil, fmt.Errorf("list notification rules: %w", err)
	}

	return rules, nil
}

func (s *Service) SendTestNotification(
	ctx context.Context,
	projectID domain.ProjectID,
	notificationSettingID domain.NotificationSettingID,
) error {
	project, err := s.projectsRepo.GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("get project by ID: %w", err)
	}

	settings, err := s.notificationSettingsRepo.GetSettingByID(ctx, notificationSettingID)
	if err != nil {
		return fmt.Errorf("get notification setting: %w", err)
	}

	for _, channel := range s.notificationChannels {
		if channel.Type() == settings.Type {
			issue := domain.Issue{
				ID:                 999,
				ProjectID:          projectID,
				Fingerprint:        "test",
				Source:             domain.SourceEvent,
				Status:             domain.IssueStatusUnresolved,
				Title:              "Test issue",
				Level:              domain.IssueLevelDebug,
				Platform:           "go",
				FirstSeen:          time.Now(),
				LastSeen:           time.Now(),
				TotalEvents:        1,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
				LastNotificationAt: nil,
			}

			return channel.Send(ctx, &issue, &project, settings.Config, false)
		}
	}

	return errors.New("channel not found")
}

package settings

import (
	"context"

	"github.com/rom8726/warden/internal/backend/contract"
	"github.com/rom8726/warden/internal/domain"
)

// Service provides settings management functionality.
type Service struct {
	settingsRepo contract.SettingRepository
	secret       []byte
}

// New creates a new settings use case.
func New(settingsRepo contract.SettingRepository, secret string) *Service {
	return &Service{
		settingsRepo: settingsRepo,
		secret:       []byte(secret),
	}
}

// GetSetting retrieves a setting by name.
func (s *Service) GetSetting(ctx context.Context, name string) (*domain.Setting, error) {
	return s.settingsRepo.GetByName(ctx, name)
}

// SetSetting creates or updates a setting.
func (s *Service) SetSetting(ctx context.Context, name string, value interface{}, description string) error {
	return s.settingsRepo.SetByName(ctx, name, value, description)
}

// DeleteSetting deletes a setting by name.
func (s *Service) DeleteSetting(ctx context.Context, name string) error {
	return s.settingsRepo.DeleteByName(ctx, name)
}

// ListSettings retrieves all settings.
func (s *Service) ListSettings(ctx context.Context) ([]*domain.Setting, error) {
	return s.settingsRepo.List(ctx)
}

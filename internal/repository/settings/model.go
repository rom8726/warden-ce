package settings

import (
	"encoding/json"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

// settingModel represents the settings table structure.
type settingModel struct {
	ID          int             `db:"id"`
	Name        string          `db:"name"`
	Value       json.RawMessage `db:"value"`
	Description string          `db:"description"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

// toDomain converts the model to a domain entity.
func (m *settingModel) toDomain() *domain.Setting {
	return &domain.Setting{
		ID:          m.ID,
		Name:        m.Name,
		Value:       m.Value,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

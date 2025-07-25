package domain

import (
	"encoding/json"
	"time"
)

// Setting represents a configuration setting stored in the database.
type Setting struct {
	ID          int             `db:"id"          json:"id"`
	Name        string          `db:"name"        json:"name"`
	Value       json.RawMessage `db:"value"       json:"value"`
	Description string          `db:"description" json:"description"`
	CreatedAt   time.Time       `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"  json:"updated_at"`
}

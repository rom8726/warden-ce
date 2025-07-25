package domain

import (
	"strconv"
	"time"
)

type ProjectID uint

// Project represents a Sentry project.
type Project struct {
	ID          ProjectID
	Name        string
	Description string
	PublicKey   string
	TeamID      *TeamID
	CreatedAt   time.Time
	ArchivedAt  *time.Time
}

type ProjectExtended struct {
	Project
	TeamName *string
}

type ProjectDTO struct {
	Name        string
	Description string
	PublicKey   string
	TeamID      *TeamID
}

func (id ProjectID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

func (id ProjectID) Uint() uint {
	return uint(id)
}

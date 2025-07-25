package domain

import "time"

type ReleaseID uint

type Release struct {
	ID          ReleaseID
	ProjectID   ProjectID
	Version     string
	Description string
	ReleasedAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReleaseDTO struct {
	ProjectID   ProjectID
	Version     string
	Description string
}

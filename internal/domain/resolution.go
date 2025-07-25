package domain

import (
	"time"
)

type ResolutionID uint

// Resolution represents a resolution of an error.
type Resolution struct {
	ID         ResolutionID
	ProjectID  ProjectID
	IssueID    IssueID
	Status     IssueStatus
	ResolvedBy *UserID
	ResolvedAt *time.Time
	Comment    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ResolutionDTO struct {
	ProjectID  ProjectID
	IssueID    IssueID
	Status     IssueStatus
	ResolvedBy *UserID
	Comment    string
}

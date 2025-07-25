package domain

import (
	"time"
)

type IssueID uint

type IssueSource string

const (
	SourceEvent     IssueSource = "event"
	SourceException IssueSource = "exception"
)

// IssueStatus represents the status of a resolution.
type IssueStatus string

const (
	IssueStatusResolved   IssueStatus = "resolved"
	IssueStatusUnresolved IssueStatus = "unresolved"
	IssueStatusIgnored    IssueStatus = "ignored"
)

type IssueLevel string

const (
	IssueLevelFatal     IssueLevel = "fatal"
	IssueLevelException IssueLevel = "exception"
	IssueLevelError     IssueLevel = "error"
	IssueLevelWarning   IssueLevel = "warning"
	IssueLevelInfo      IssueLevel = "info"
	IssueLevelDebug     IssueLevel = "debug"
)

type Issue struct {
	ID                 IssueID
	ProjectID          ProjectID
	Fingerprint        string
	Source             IssueSource
	Status             IssueStatus
	Title              string
	Level              IssueLevel
	Platform           string
	FirstSeen          time.Time
	LastSeen           time.Time
	TotalEvents        uint
	CreatedAt          time.Time
	UpdatedAt          time.Time
	LastNotificationAt *time.Time
}

type IssueExtended struct {
	Issue
	ProjectName        string
	ResolvedBy         *UserID
	ResolvedByUsername *string
	ResolvedAt         *time.Time
}

type IssueDTO struct {
	ProjectID   ProjectID
	Fingerprint string
	Source      IssueSource
	Status      IssueStatus
	Title       string
	Level       IssueLevel
	Platform    string
}

type IssueExtendedWithChildren struct {
	Issue
	ProjectName        string
	ResolvedBy         *UserID
	ResolvedByUsername *string
	ResolvedAt         *time.Time
	Events             []Event
}

func (id IssueID) Uint() uint {
	return uint(id)
}

type IssueUpsertResult struct {
	ID             IssueID
	IsNew          bool
	WasReactivated bool
}

func IsNotifiableLevel(level IssueLevel) bool {
	if level == IssueLevelDebug || level == IssueLevelInfo {
		return false
	}

	return true
}

func (lvl IssueLevel) String() string {
	return string(lvl)
}

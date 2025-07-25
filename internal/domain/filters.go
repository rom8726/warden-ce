package domain

import (
	"time"
)

type OrderByColumn uint8

const (
	OrderByFieldTotalEvents OrderByColumn = iota // total_events column
	OrderByFieldFirstSeen                        // first_seen column
	OrderByFieldLastSeen                         // last_seen column
)

type ListIssuesFilter struct {
	ProjectID *ProjectID
	Level     *IssueLevel
	Status    *IssueStatus
	TimeFrom  time.Time
	TimeTo    time.Time
	OrderBy   OrderByColumn // total_events by default
	OrderAsc  bool          // ASC / DESC, DESC by default
	PageNum   uint          // from 1
	PerPage   uint          // records limit
}

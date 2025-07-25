package domain

import (
	"time"
)

type Period struct {
	Interval    time.Duration
	Granularity time.Duration
}

type Timeseries struct {
	Name        string
	Period      Period
	Occurrences []uint
}

type IssueTimeseriesFilter struct {
	Period    Period
	ProjectID *ProjectID
	IssueID   *IssueID
	Levels    []IssueLevel
	Statuses  []IssueStatus

	GroupBy IssueTimeseriesGroup
}

type IssueTimeseriesGroup uint8

const (
	IssueTimeseriesGroupNone IssueTimeseriesGroup = iota
	IssueTimeseriesGroupProject
	IssueTimeseriesGroupIssue
	IssueTimeseriesGroupLevel
	IssueTimeseriesGroupStatus
)

type EventTimeseriesFilter struct {
	Period Period

	ProjectID *ProjectID
	Levels    []IssueLevel
	Release   *string

	GroupBy EventTimeseriesGroup
}

type EventTimeseriesGroup uint8

const (
	EventTimeseriesGroupNone EventTimeseriesGroup = iota
	EventTimeseriesGroupProject
	EventTimeseriesGroupLevel
)

type IssueEventsTimeseriesFilter struct {
	Period    Period
	ProjectID ProjectID
	IssueID   IssueID
	Levels    []IssueLevel

	GroupBy EventTimeseriesGroup
}

type IssueEventsTimeseriesGroup uint8

const (
	IssueEventsTimeseriesGroupNone IssueEventsTimeseriesGroup = iota
	IssueEventsTimeseriesGroupLevel
)

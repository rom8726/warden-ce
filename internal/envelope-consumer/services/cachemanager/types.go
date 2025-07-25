package cachemanager

import (
	"fmt"
)

// Type represents the type of cache.
type Type string

const (
	TypeRelease      Type = "release"
	TypeIssue        Type = "issue"
	TypeIssueRelease Type = "issue_release"
)

// ReleaseKey represents a key for the release cache.
type ReleaseKey struct {
	ProjectID uint
	Version   string
}

func (k ReleaseKey) String() string {
	return fmt.Sprintf("%d:%s", k.ProjectID, k.Version)
}

// IssueKey represents a key for the issue cache.
type IssueKey struct {
	Fingerprint string
}

func (k IssueKey) String() string {
	return k.Fingerprint
}

// IssueReleaseKey represents a key for issue_release cache.
type IssueReleaseKey struct {
	IssueID   uint
	ReleaseID uint
}

func (k IssueReleaseKey) String() string {
	return fmt.Sprintf("%d:%d", k.IssueID, k.ReleaseID)
}

// ReleaseValue represents a value for release cache.
type ReleaseValue struct {
	ReleaseID uint
}

func (v ReleaseValue) IsValid() bool {
	return v.ReleaseID > 0
}

// IssueValue represents a value for issue cache.
type IssueValue struct {
	IssueID uint
}

func (v IssueValue) IsValid() bool {
	return v.IssueID > 0
}

// IssueReleaseValue represents a value for issue_release cache.
type IssueReleaseValue struct {
	IssueReleaseID uint
	FirstSeenIn    bool
}

func (v IssueReleaseValue) IsValid() bool {
	return v.IssueReleaseID > 0
}

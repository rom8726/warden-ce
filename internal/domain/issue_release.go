package domain

type IssueRelease struct {
	IssueID     IssueID
	ReleaseID   ReleaseID
	FirstSeenIn bool
}

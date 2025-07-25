package issuereleases

// nolint:unused // will be used
type issueReleaseModel struct {
	IssueID     uint `db:"issue_id"`
	ReleaseID   uint `db:"release_id"`
	FirstSeenIn bool `db:"first_seen_in"`
}

package domain

type GeneralProjectStats struct {
	TotalIssues        uint
	FatalIssues        uint
	ErrorIssues        uint
	WarningIssues      uint
	InfoIssues         uint
	DebugIssues        uint
	ExceptionIssues    uint
	MostFrequentIssues []IssueExtended
}

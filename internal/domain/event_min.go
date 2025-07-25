package domain

type EventMinimal struct {
	Level    IssueLevel
	Source   IssueSource
	Platform string
	Message  string

	ExceptionData
}

func (ev *EventMinimal) FullFingerprint() string {
	eventFull := Event{
		Level:         ev.Level,
		Source:        ev.Source,
		Platform:      ev.Platform,
		Message:       ev.Message,
		ExceptionData: ev.ExceptionData,
	}

	return eventFull.FullFingerprint()
}

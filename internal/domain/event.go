package domain

import (
	"encoding/json"
	"time"

	"github.com/rom8726/warden/pkg/fingerprinter"
)

type EventID string

// Event represents a Sentry event.
type Event struct {
	Timestamp time.Time
	GroupHash string

	ID          EventID
	ProjectID   ProjectID
	Level       IssueLevel
	Source      IssueSource
	Platform    string
	Message     string
	Payload     json.RawMessage
	Tags        map[string]string
	ServerName  string
	Environment string
	Release     string

	ExceptionData
	EventRequestContext
	EventUserData
	EventRuntimeContext
}

type ExceptionData struct {
	ExceptionType       *string
	ExceptionValue      *string
	ExceptionStacktrace json.RawMessage
}

type EventRequestContext struct {
	RequestURL     *string
	RequestMethod  *string
	RequestQuery   *string
	RequestHeaders map[string]string
	RequestData    *string
	RequestCookies *string
	RequestIP      *string
	UserAgent      *string
}

type EventUserData struct {
	UserID    *string
	UserEmail *string
}

type EventRuntimeContext struct {
	RuntimeName    *string
	RuntimeVersion *string
	OSName         *string
	OSVersion      *string
	BrowserName    *string
	BrowserVersion *string
	DeviceArch     *string
}

func (ev *Event) FullFingerprint() string {
	if ev.Source == SourceEvent {
		return fingerprinter.SHA1FromStrings(ev.Message, string(ev.Level), ev.Platform)
	}

	stacktraceData, _ := ev.ExceptionStacktrace.MarshalJSON()

	return fingerprinter.SHA1FromStrings(*ev.ExceptionType, *ev.ExceptionValue, string(stacktraceData))
}

func (id EventID) String() string {
	return string(id)
}

//nolint:gocritic // need refactor
package events

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type eventModel struct {
	EventID   string    `ch:"event_id"   json:"event_id"`
	ProjectID uint32    `ch:"project_id" json:"project_id"`
	Message   string    `ch:"message"    json:"message"`
	Level     string    `ch:"level"      json:"level"`
	Platform  string    `ch:"platform"   json:"platform"`
	Timestamp time.Time `ch:"timestamp"  json:"timestamp"`
	GroupHash string    `ch:"group_hash" json:"group_hash"`
	Source    string    `ch:"source"     json:"source"` // 'event' or 'exception'

	// Exception context (nullable)
	Stacktrace     *string `ch:"stacktrace"      json:"stacktrace,omitempty"`
	ExceptionType  *string `ch:"exception_type"  json:"exception_type,omitempty"`
	ExceptionValue *string `ch:"exception_value" json:"exception_value,omitempty"`

	// Request context
	RequestURL     *string           `ch:"request_url"     json:"request_url,omitempty"`
	RequestMethod  *string           `ch:"request_method"  json:"request_method,omitempty"`
	RequestQuery   *string           `ch:"request_query"   json:"request_query,omitempty"`
	RequestHeaders map[string]string `ch:"request_headers" json:"request_headers,omitempty"`
	RequestData    *string           `ch:"request_data"    json:"request_data,omitempty"`
	RequestCookies *string           `ch:"request_cookies" json:"request_cookies,omitempty"`
	RequestIP      *string           `ch:"request_ip"      json:"request_ip,omitempty"`

	// User
	UserID    *string `ch:"user_id"    json:"user_id,omitempty"`
	UserEmail *string `ch:"user_email" json:"user_email,omitempty"`
	UserAgent *string `ch:"user_agent" json:"user_agent,omitempty"`

	// Contexts
	RuntimeName    *string `ch:"runtime_name"    json:"runtime_name,omitempty"`
	RuntimeVersion *string `ch:"runtime_version" json:"runtime_version,omitempty"`
	OSName         *string `ch:"os_name"         json:"os_name,omitempty"`
	OSVersion      *string `ch:"os_version"      json:"os_version,omitempty"`
	BrowserName    *string `ch:"browser_name"    json:"browser_name,omitempty"`
	BrowserVersion *string `ch:"browser_version" json:"browser_version,omitempty"`
	DeviceArch     *string `ch:"device_arch"     json:"device_arch,omitempty"`

	// Common
	ServerName  string            `ch:"server_name" json:"server_name"`
	Environment string            `ch:"environment" json:"environment"`
	Tags        map[string]string `ch:"tags"        json:"tags"`

	// Raw JSON
	RawData string `ch:"raw_data" json:"raw_data"`

	// App
	Release *string `ch:"release" json:"release"`
}

func (m eventModel) MarshalJSON() ([]byte, error) {
	type Alias eventModel

	return json.Marshal(&struct {
		Timestamp string `json:"timestamp"`
		Alias
	}{
		Timestamp: m.Timestamp.UTC().Format("2006-01-02 15:04:05"),
		Alias:     Alias(m),
	})
}

func (m *eventModel) toDomain() (domain.Event, error) {
	var exception domain.ExceptionData
	if domain.IssueSource(m.Source) == domain.SourceException {
		stacktrace, err := gzip64ToJsonRawMessage(*m.Stacktrace)
		if err != nil {
			return domain.Event{}, err
		}

		exception = domain.ExceptionData{
			ExceptionType:       m.ExceptionType,
			ExceptionValue:      m.ExceptionValue,
			ExceptionStacktrace: stacktrace,
		}
	}

	var release string
	if m.Release != nil {
		release = *m.Release
	}

	//payload, err := gzip64ToJsonRawMessage(m.RawData)
	//if err != nil {
	//	panic(err)
	//}
	var payload json.RawMessage

	return domain.Event{
		Timestamp:     m.Timestamp,
		GroupHash:     m.GroupHash,
		ID:            domain.EventID(m.EventID),
		ProjectID:     domain.ProjectID(m.ProjectID),
		Level:         domain.IssueLevel(m.Level),
		Source:        domain.IssueSource(m.Source),
		Platform:      m.Platform,
		Release:       release,
		Message:       m.Message,
		Payload:       payload,
		Tags:          m.Tags,
		ServerName:    m.ServerName,
		Environment:   m.Environment,
		ExceptionData: exception,
		EventRequestContext: domain.EventRequestContext{
			RequestURL:     m.RequestURL,
			RequestMethod:  m.RequestMethod,
			RequestQuery:   m.RequestQuery,
			RequestHeaders: m.RequestHeaders,
			RequestData:    m.RequestData,
			RequestCookies: m.RequestCookies,
			RequestIP:      m.RequestIP,
			UserAgent:      m.UserAgent,
		},
		EventUserData: domain.EventUserData{
			UserID:    m.UserID,
			UserEmail: m.UserEmail,
		},
		EventRuntimeContext: domain.EventRuntimeContext{
			RuntimeName:    m.RuntimeName,
			RuntimeVersion: m.RuntimeVersion,
			OSName:         m.OSName,
			OSVersion:      m.OSVersion,
			BrowserName:    m.BrowserName,
			BrowserVersion: m.BrowserVersion,
			DeviceArch:     m.DeviceArch,
		},
	}, nil
}

func fromDomain(event *domain.Event) (eventModel, error) {
	var stacktrace *string
	if event.ExceptionStacktrace != nil {
		stacktraceStr, err := jsonRawMessageToGzip64(event.ExceptionStacktrace)
		if err != nil {
			return eventModel{}, err
		}

		stacktrace = &stacktraceStr
	}

	var release *string
	if event.Release != "" {
		release = &event.Release
	}

	//rawData, err := jsonRawMessageToGzip64(event.Payload)
	//if err != nil {
	//	return eventModel{}, err
	//}
	var rawData string

	return eventModel{
		EventID:        string(event.ID),
		ProjectID:      uint32(event.ProjectID), //nolint:gosec // it's ok
		Message:        event.Message,
		Level:          string(event.Level),
		Platform:       event.Platform,
		Timestamp:      event.Timestamp,
		GroupHash:      event.GroupHash,
		Source:         string(event.Source),
		Stacktrace:     stacktrace,
		ExceptionType:  event.ExceptionType,
		ExceptionValue: event.ExceptionValue,
		RequestURL:     event.RequestURL,
		RequestMethod:  event.RequestMethod,
		RequestQuery:   event.RequestQuery,
		RequestHeaders: event.RequestHeaders,
		RequestData:    event.RequestData,
		RequestCookies: event.RequestCookies,
		RequestIP:      event.RequestIP,
		UserID:         event.UserID,
		UserEmail:      event.UserEmail,
		UserAgent:      event.UserAgent,
		RuntimeName:    event.RuntimeName,
		RuntimeVersion: event.RuntimeVersion,
		OSName:         event.OSName,
		OSVersion:      event.OSVersion,
		BrowserName:    event.BrowserName,
		BrowserVersion: event.BrowserVersion,
		DeviceArch:     event.DeviceArch,
		ServerName:     event.ServerName,
		Environment:    event.Environment,
		Tags:           event.Tags,
		RawData:        rawData,
		Release:        release,
	}, nil
}

func jsonRawMessageToGzip64(data json.RawMessage) (string, error) {
	var buf bytes.Buffer

	gzipWriter := gzip.NewWriter(&buf)

	_, err := gzipWriter.Write(data)
	if err != nil {
		return "", err
	}

	err = gzipWriter.Close()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func gzip64ToJsonRawMessage(encoded string) (json.RawMessage, error) {
	compressedData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()

	var uncompressedData bytes.Buffer
	_, err = io.Copy(&uncompressedData, reader) //nolint:gosec //it's ok
	if err != nil {
		return nil, err
	}

	return uncompressedData.Bytes(), nil
}

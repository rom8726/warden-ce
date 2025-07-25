package dto

import (
	"encoding/json"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

// DomainIssueEventToAPI converts domain.Event to generatedapi.IssueEvent.
func DomainIssueEventToAPI(event domain.Event) generatedapi.IssueEvent {
	// --- Tags ---
	var optTags generatedapi.OptIssueEventTags
	if event.Tags != nil {
		tags := make(generatedapi.IssueEventTags, len(event.Tags))
		for k, v := range event.Tags {
			tags[k] = v
		}

		optTags.Value = tags
		optTags.Set = true
	}

	// --- Exception ---
	var (
		excType generatedapi.OptNilString
		excVal  generatedapi.OptNilString
		excStk  generatedapi.OptNilString
	)
	if event.ExceptionData.ExceptionType != nil {
		excType.Value = *event.ExceptionData.ExceptionType
		excType.Set = true

		excVal.Value = *event.ExceptionData.ExceptionValue
		excVal.Set = true

		excStk.Value = string(event.ExceptionData.ExceptionStacktrace)
		excStk.Set = true
	}

	// --- Simple optionals ---
	makeOptString := func(v string) (opt generatedapi.OptString) {
		if v != "" {
			opt.Value = v
			opt.Set = true
		}

		return
	}

	optServerName := makeOptString(event.ServerName)
	optEnvironment := makeOptString(event.Environment)
	optRelease := makeOptString(event.Release)
	optGroupHash := makeOptString(event.GroupHash)

	// --- Payload ---
	var optPayload generatedapi.OptIssueEventPayload
	if len(event.Payload) != 0 {
		var p generatedapi.IssueEventPayload
		_ = json.Unmarshal(event.Payload, &p) // best-effort
		optPayload.Value = p
		optPayload.Set = true
	}

	// --- Request context ---
	makeOptNilString := func(ptr *string) (opt generatedapi.OptNilString) {
		if ptr != nil {
			opt.Value = *ptr
			opt.Set = true
		}

		return
	}

	reqURL := makeOptNilString(event.RequestURL)
	reqMethod := makeOptNilString(event.RequestMethod)
	reqQuery := makeOptNilString(event.RequestQuery)
	reqData := makeOptNilString(event.RequestData)
	reqCookies := makeOptNilString(event.RequestCookies)
	reqIP := makeOptNilString(event.RequestIP)
	userAgent := makeOptNilString(event.UserAgent)

	var reqHeaders generatedapi.OptNilIssueEventRequestHeaders
	if event.RequestHeaders != nil {
		h := make(generatedapi.IssueEventRequestHeaders, len(event.RequestHeaders))
		for k, v := range event.RequestHeaders {
			h[k] = v
		}
		reqHeaders.Value = h
		reqHeaders.Set = true
	}

	// --- User data ---
	userID := makeOptNilString(event.UserID)
	userEmail := makeOptNilString(event.UserEmail)

	// --- Runtime context ---
	runtimeName := makeOptNilString(event.RuntimeName)
	runtimeVersion := makeOptNilString(event.RuntimeVersion)
	osName := makeOptNilString(event.OSName)
	osVersion := makeOptNilString(event.OSVersion)
	browserName := makeOptNilString(event.BrowserName)
	browserVersion := makeOptNilString(event.BrowserVersion)
	deviceArch := makeOptNilString(event.DeviceArch)

	return generatedapi.IssueEvent{
		EventID:     event.ID.String(),
		GroupHash:   optGroupHash,
		ProjectID:   event.ProjectID.Uint(),
		Message:     event.Message,
		Level:       generatedapi.IssueLevel(event.Level),
		Source:      generatedapi.IssueSource(event.Source),
		Platform:    event.Platform,
		Payload:     optPayload,
		Timestamp:   event.Timestamp,
		ServerName:  optServerName,
		Environment: optEnvironment,
		Tags:        optTags,
		Release:     optRelease,

		ExceptionType:       excType,
		ExceptionValue:      excVal,
		ExceptionStacktrace: excStk,

		RequestURL:     reqURL,
		RequestMethod:  reqMethod,
		RequestQuery:   reqQuery,
		RequestHeaders: reqHeaders,
		RequestData:    reqData,
		RequestCookies: reqCookies,
		RequestIP:      reqIP,
		UserAgent:      userAgent,

		UserID:         userID,
		UserEmail:      userEmail,
		RuntimeName:    runtimeName,
		RuntimeVersion: runtimeVersion,
		OsName:         osName,
		OsVersion:      osVersion,
		BrowserName:    browserName,
		BrowserVersion: browserVersion,
		DeviceArch:     deviceArch,
	}
}

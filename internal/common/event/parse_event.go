package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

const unknownKeyword = "unknown"

func ParseEvent(eventData map[string]any, projectID domain.ProjectID) (domain.Event, error) {
	// Extract event ID
	eventIDRaw, ok := eventData["event_id"]
	if !ok {
		return domain.Event{}, errors.New("event_id is required")
	}
	eventID := domain.EventID(fmt.Sprint(eventIDRaw))

	// Extract message
	message := extractMessage(eventData)

	// Extract level (default to "error" if not present)
	level := extractLevel(eventData)

	// Extract platform (default to "unknown" if not present)
	platform := extractPlatform(eventData)

	// Extract timestamp (default now if not present)
	timestamp := extractTimestamp(eventData)

	// Extract tags if available
	tags := extractTags(eventData)

	// Extract environment
	environment := extractEnvironment(eventData)

	// Extract server name
	serverName := extractServerName(eventData)

	// Extract release
	release := extractRelease(eventData)

	// Convert the event data to JSON
	rawData, err := json.Marshal(eventData)
	if err != nil {
		return domain.Event{}, fmt.Errorf("marshal event data: %w", err)
	}

	source := domain.SourceEvent
	var exceptionData domain.ExceptionData
	if _, ok := eventData["exception"]; ok { //nolint:nestif // need refactoring
		if level != domain.IssueLevelFatal {
			source = domain.SourceException
			level = domain.IssueLevelException
		}

		exValues, err := extractExceptionValues(eventData)
		if err != nil {
			return domain.Event{}, fmt.Errorf("extract exception values: %w", err)
		}
		if len(exValues) > 0 {
			value := exValues[0]
			exceptionData.ExceptionStacktrace = value.Stacktrace
			exceptionData.ExceptionType = &value.Type
			exceptionData.ExceptionValue = &value.Value
			if message == "" {
				message = value.Value
			}
		}
	}

	// ------------------------------------------------------------------
	// Request context
	// ------------------------------------------------------------------
	reqCtx := extractRequestContext(eventData)

	// ------------------------------------------------------------------
	// User data
	// ------------------------------------------------------------------
	userCtx := extractUserContext(eventData, &reqCtx)

	// ------------------------------------------------------------------
	// Runtime / OS / Browser / Device
	// ------------------------------------------------------------------
	runtimeCtx := extractRuntimeContext(eventData)

	// Create the event
	event := domain.Event{
		Timestamp:           timestamp,
		ID:                  eventID,
		ProjectID:           projectID,
		Level:               level,
		Source:              source,
		Platform:            platform,
		Message:             message,
		Payload:             rawData,
		Tags:                tags,
		ServerName:          serverName,
		Release:             release,
		Environment:         environment,
		ExceptionData:       exceptionData,
		EventRequestContext: reqCtx,
		EventUserData:       userCtx,
		EventRuntimeContext: runtimeCtx,
	}

	event.GroupHash = event.FullFingerprint()

	return event, nil
}

type ExceptionValue struct {
	Type       string
	Value      string
	Stacktrace json.RawMessage
}

// extractExceptionValues extracts exception values from the event data.
//
//nolint:nestif // need refactoring
func extractExceptionValues(data map[string]any) ([]ExceptionValue, error) {
	var result []ExceptionValue

	// Check if there's an "exception" field
	exceptionsRaw, ok := data["exception"]
	if !ok {
		return result, nil
	}

	if exceptionsMap, ok := exceptionsRaw.(map[string]any); ok {
		return extractFromExceptionMap(exceptionsMap)
	} else if exceptionsList, ok := exceptionsRaw.([]any); ok {
		return extractFromExceptionList(exceptionsList)
	}

	return result, nil
}

func extractFromExceptionMap(exceptionsMap map[string]any) ([]ExceptionValue, error) {
	var result []ExceptionValue
	if valuesRaw, ok := exceptionsMap["values"]; ok {
		if valuesList, ok := valuesRaw.([]any); ok {
			return extractFromExceptionValuesList(valuesList)
		}
	} else {
		// Single exception
		exValue, err := extractSingleException(exceptionsMap)
		if err != nil {
			return nil, err
		}
		result = append(result, exValue)
	}

	return result, nil
}

func extractFromExceptionValuesList(valuesList []any) ([]ExceptionValue, error) {
	var result []ExceptionValue
	for _, valRaw := range valuesList {
		if valMap, ok := valRaw.(map[string]any); ok {
			exValue, err := extractSingleException(valMap)
			if err != nil {
				return nil, err
			}
			result = append(result, exValue)
		}
	}

	return result, nil
}

// extractSingleException extracts a single exception from a map.
func extractSingleException(exMap map[string]any) (ExceptionValue, error) {
	var result ExceptionValue

	// Extract type
	if typeRaw, ok := exMap["type"]; ok {
		result.Type = fmt.Sprint(typeRaw)
	} else {
		return result, errors.New("exception type is required")
	}

	// Extract value
	if valueRaw, ok := exMap["value"]; ok {
		result.Value = fmt.Sprint(valueRaw)
	} else {
		return result, errors.New("exception value is required")
	}

	// Extract stacktrace if available
	if stacktraceRaw, ok := exMap["stacktrace"]; ok {
		stacktraceJSON, err := json.Marshal(stacktraceRaw)
		if err != nil {
			return result, fmt.Errorf("marshal stacktrace: %w", err)
		}
		result.Stacktrace = stacktraceJSON
	}

	return result, nil
}

func extractRequestContext(eventData map[string]any) domain.EventRequestContext {
	var reqCtx domain.EventRequestContext
	if reqRaw, ok := eventData["request"].(map[string]any); ok {
		extractRequestURL(reqRaw, &reqCtx)
		extractRequestMethod(reqRaw, &reqCtx)
		extractRequestQuery(reqRaw, &reqCtx)
		extractRequestData(reqRaw, &reqCtx)
		extractRequestCookies(reqRaw, &reqCtx)
		extractRequestHeaders(reqRaw, &reqCtx)
		extractRequestEnv(reqRaw, &reqCtx)
		extractRequestIP(reqRaw, &reqCtx)
	}

	return reqCtx
}

func extractRequestURL(reqRaw map[string]any, reqCtx *domain.EventRequestContext) {
	if v, ok := reqRaw["url"].(string); ok {
		reqCtx.RequestURL = &v
	}
}

func extractRequestMethod(reqRaw map[string]any, reqCtx *domain.EventRequestContext) {
	if v, ok := reqRaw["method"].(string); ok {
		reqCtx.RequestMethod = &v
	}
}

func extractRequestQuery(reqRaw map[string]any, reqCtx *domain.EventRequestContext) {
	if v, ok := reqRaw["query_string"].(string); ok {
		reqCtx.RequestQuery = &v
	}
}

func extractRequestData(reqRaw map[string]any, reqCtx *domain.EventRequestContext) {
	if v, ok := reqRaw["data"].(string); ok {
		reqCtx.RequestData = &v
	}
}

func extractRequestCookies(reqRaw map[string]any, reqCtx *domain.EventRequestContext) {
	if v, ok := reqRaw["cookies"].(string); ok {
		reqCtx.RequestCookies = &v
	}
}

func extractRequestHeaders(reqRaw map[string]any, reqCtx *domain.EventRequestContext) {
	if hdrRaw, ok := reqRaw["headers"].(map[string]any); ok {
		h := make(map[string]string, len(hdrRaw))
		for k, v := range hdrRaw {
			h[k] = fmt.Sprint(v)
		}
		reqCtx.RequestHeaders = h
		if ua, ok := h["User-Agent"]; ok {
			reqCtx.UserAgent = &ua
		}
	}
}

func extractRequestEnv(reqRaw map[string]any, reqCtx *domain.EventRequestContext) {
	if envRaw, ok := reqRaw["env"].(map[string]any); ok {
		if ip, ok := envRaw["REMOTE_ADDR"].(string); ok {
			reqCtx.RequestIP = &ip
		}
	}
}

func extractRequestIP(reqRaw map[string]any, reqCtx *domain.EventRequestContext) {
	if reqCtx.RequestIP == nil {
		if ip, ok := reqRaw["ip_address"].(string); ok {
			reqCtx.RequestIP = &ip
		}
	}
}

func extractUserContext(eventData map[string]any, reqCtx *domain.EventRequestContext) domain.EventUserData {
	var userCtx domain.EventUserData
	if userRaw, ok := eventData["user"].(map[string]any); ok {
		extractUserID(userRaw, &userCtx)
		extractUserEmail(userRaw, &userCtx)
		extractUserIP(userRaw, reqCtx)
	}

	return userCtx
}

func extractUserID(userRaw map[string]any, userCtx *domain.EventUserData) {
	if v, ok := userRaw["id"].(string); ok {
		userCtx.UserID = &v
	}
}

func extractUserEmail(userRaw map[string]any, userCtx *domain.EventUserData) {
	if v, ok := userRaw["email"].(string); ok {
		userCtx.UserEmail = &v
	}
}

func extractUserIP(userRaw map[string]any, reqCtx *domain.EventRequestContext) {
	if reqCtx.RequestIP == nil {
		if ip, ok := userRaw["ip_address"].(string); ok {
			reqCtx.RequestIP = &ip
		}
	}
}

func extractRuntimeContext(eventData map[string]any) domain.EventRuntimeContext {
	var runtimeCtx domain.EventRuntimeContext
	if ctxRaw, ok := eventData["contexts"].(map[string]any); ok {
		extractRuntime(ctxRaw, &runtimeCtx)
		extractOS(ctxRaw, &runtimeCtx)
		extractBrowser(ctxRaw, &runtimeCtx)
		extractDevice(ctxRaw, &runtimeCtx)
	}

	return runtimeCtx
}

func extractRuntime(ctxRaw map[string]any, runtimeCtx *domain.EventRuntimeContext) {
	if rt, ok := ctxRaw["runtime"].(map[string]any); ok {
		if v, ok := rt["name"].(string); ok {
			runtimeCtx.RuntimeName = &v
		}
		if v, ok := rt["version"].(string); ok {
			runtimeCtx.RuntimeVersion = &v
		}
	}
}

func extractOS(ctxRaw map[string]any, runtimeCtx *domain.EventRuntimeContext) {
	if os, ok := ctxRaw["os"].(map[string]any); ok {
		if v, ok := os["name"].(string); ok {
			runtimeCtx.OSName = &v
		}
		if v, ok := os["version"].(string); ok {
			runtimeCtx.OSVersion = &v
		}
	}
}

func extractBrowser(ctxRaw map[string]any, runtimeCtx *domain.EventRuntimeContext) {
	if br, ok := ctxRaw["browser"].(map[string]any); ok {
		if v, ok := br["name"].(string); ok {
			runtimeCtx.BrowserName = &v
		}
		if v, ok := br["version"].(string); ok {
			runtimeCtx.BrowserVersion = &v
		}
	}
}

func extractDevice(ctxRaw map[string]any, runtimeCtx *domain.EventRuntimeContext) {
	if dv, ok := ctxRaw["device"].(map[string]any); ok {
		if v, ok := dv["arch"].(string); ok {
			runtimeCtx.DeviceArch = &v
		}
	}
}

func extractFromExceptionList(exceptionsList []any) ([]ExceptionValue, error) {
	return extractFromExceptionValuesList(exceptionsList)
}

func extractTags(eventData map[string]any) map[string]string {
	tags := make(map[string]string)
	if tagsRaw, ok := eventData["tags"]; ok {
		if tagsMap, ok := tagsRaw.(map[string]any); ok {
			for k, v := range tagsMap {
				tags[k] = fmt.Sprint(v)
			}
		}
	}

	return tags
}

func extractTimestamp(eventData map[string]any) time.Time {
	timestamp := time.Now()
	if timestampRaw, ok := eventData["timestamp"]; ok {
		if ts, ok := timestampRaw.(string); ok {
			if parsedTime, err := time.Parse(time.RFC3339, ts); err == nil {
				timestamp = parsedTime
			}
		}
	}

	return timestamp
}

func extractMessage(eventData map[string]any) string {
	if messageRaw, ok := eventData["message"]; ok {
		return fmt.Sprint(messageRaw)
	}

	return ""
}

func extractLevel(eventData map[string]any) domain.IssueLevel {
	if levelRaw, ok := eventData["level"]; ok {
		return domain.IssueLevel(fmt.Sprint(levelRaw))
	}

	return "error"
}

func extractPlatform(eventData map[string]any) string {
	if platformRaw, ok := eventData["platform"]; ok {
		return fmt.Sprint(platformRaw)
	}

	return unknownKeyword
}

func extractEnvironment(eventData map[string]any) string {
	if environmentRaw, ok := eventData["environment"]; ok {
		return fmt.Sprint(environmentRaw)
	}

	return unknownKeyword
}

func extractServerName(eventData map[string]any) string {
	if serverNameRaw, ok := eventData["server_name"]; ok {
		return fmt.Sprint(serverNameRaw)
	}

	return unknownKeyword
}

func extractRelease(eventData map[string]any) string {
	if releaseRaw, ok := eventData["release"]; ok {
		return strings.TrimSpace(fmt.Sprint(releaseRaw))
	}

	return unknownKeyword
}

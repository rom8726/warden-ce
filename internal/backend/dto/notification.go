package dto

import (
	"encoding/json"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

// DomainNotificationSettingToAPI converts domain.NotificationSetting to generatedapi.NotificationSetting.
func DomainNotificationSettingToAPI(setting domain.NotificationSetting) generatedapi.NotificationSetting {
	return generatedapi.NotificationSetting{
		ID:        uint(setting.ID),
		ProjectID: uint(setting.ProjectID),
		Type:      string(setting.Type),
		Config:    string(setting.Config),
		Enabled:   setting.Enabled,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}

// DomainNotificationRuleToAPI converts domain.NotificationRule to generatedapi.NotificationRule.
func DomainNotificationRuleToAPI(rule domain.NotificationRule) generatedapi.NotificationRule {
	// Create OptNilString and OptNilBool values
	var fingerprint generatedapi.OptNilString
	if rule.Fingerprint != nil {
		fingerprint.Value = *rule.Fingerprint
		fingerprint.Set = true
	}

	var isNewError generatedapi.OptNilBool
	if rule.IsNewError != nil {
		isNewError.Value = *rule.IsNewError
		isNewError.Set = true
	}

	var isRegression generatedapi.OptNilBool
	if rule.IsRegression != nil {
		isRegression.Value = *rule.IsRegression
		isRegression.Set = true
	}

	return generatedapi.NotificationRule{
		ID:                    uint(rule.ID),
		NotificationSettingID: uint(rule.NotificationSetting),
		EventLevel:            generatedapi.NewOptNilString(string(rule.EventLevel)),
		Fingerprint:           fingerprint,
		IsNewError:            isNewError,
		IsRegression:          isRegression,
		CreatedAt:             rule.CreatedAt,
	}
}

// MakeNotificationSettingDTO converts generatedapi.CreateNotificationSettingRequest to domain.NotificationSettingDTO.
func MakeNotificationSettingDTO(
	req *generatedapi.CreateNotificationSettingRequest,
	projectID domain.ProjectID,
) domain.NotificationSettingDTO {
	return domain.NotificationSettingDTO{
		ProjectID: projectID,
		Type:      domain.NotificationType(req.Type),
		Config:    json.RawMessage(req.Config),
		Enabled:   req.Enabled.Value,
	}
}

// MakeNotificationRuleDTO converts generatedapi.CreateNotificationRuleRequest to domain.NotificationRuleDTO.
func MakeNotificationRuleDTO(
	req *generatedapi.CreateNotificationRuleRequest,
	settingID domain.NotificationSettingID,
) domain.NotificationRuleDTO {
	// Get values from OptNilString and OptNilBool
	eventLevel := ""
	if req.EventLevel.IsSet() && !req.EventLevel.IsNull() {
		eventLevel = req.EventLevel.Value
	}

	var fingerprint *string
	if req.Fingerprint.IsSet() && !req.Fingerprint.IsNull() {
		fingerprint = &req.Fingerprint.Value
	}

	var isNewError *bool
	if req.IsNewError.IsSet() && !req.IsNewError.IsNull() {
		isNewError = &req.IsNewError.Value
	}

	var isRegression *bool
	if req.IsRegression.IsSet() && !req.IsRegression.IsNull() {
		isRegression = &req.IsRegression.Value
	}

	return domain.NotificationRuleDTO{
		NotificationSetting: settingID,
		EventLevel:          eventLevel,
		Fingerprint:         fingerprint,
		IsNewError:          isNewError,
		IsRegression:        isRegression,
	}
}

// UpdateNotificationSettingFromRequest updates a domain.NotificationSetting
// from generatedapi.UpdateNotificationSettingRequest.
func UpdateNotificationSettingFromRequest(
	setting domain.NotificationSetting,
	req *generatedapi.UpdateNotificationSettingRequest,
) domain.NotificationSetting {
	if req.Type.IsSet() {
		setting.Type = domain.NotificationType(req.Type.Value)
	}

	if req.Enabled.Set {
		setting.Enabled = req.Enabled.Value
	}

	if req.Config.Set {
		setting.Config = json.RawMessage(req.Config.Value)
	}

	return setting
}

// UpdateNotificationRuleFromRequest updates a domain.NotificationRule from generatedapi.UpdateNotificationRuleRequest.
func UpdateNotificationRuleFromRequest(
	rule domain.NotificationRule,
	req *generatedapi.UpdateNotificationRuleRequest,
) domain.NotificationRule {
	if req.EventLevel.IsSet() && !req.EventLevel.IsNull() {
		rule.EventLevel = domain.IssueLevel(req.EventLevel.Value)
	}

	if req.Fingerprint.IsSet() && !req.Fingerprint.IsNull() {
		rule.Fingerprint = &req.Fingerprint.Value
	}

	if req.IsNewError.IsSet() && !req.IsNewError.IsNull() {
		rule.IsNewError = &req.IsNewError.Value
	}

	if req.IsRegression.IsSet() && !req.IsRegression.IsNull() {
		rule.IsRegression = &req.IsRegression.Value
	}

	return rule
}

// MakeListNotificationSettingsResponse converts a slice of domain.NotificationSetting
// to generatedapi.ListNotificationSettingsResponse.
func MakeListNotificationSettingsResponse(
	settings []domain.NotificationSetting,
) generatedapi.ListNotificationSettingsResponse {
	apiSettings := make([]generatedapi.NotificationSetting, len(settings))
	for i, setting := range settings {
		apiSettings[i] = DomainNotificationSettingToAPI(setting)
	}

	return generatedapi.ListNotificationSettingsResponse{
		NotificationSettings: apiSettings,
	}
}

// MakeListNotificationRulesResponse converts a slice of domain.NotificationRule
// to generatedapi.ListNotificationRulesResponse.
func MakeListNotificationRulesResponse(rules []domain.NotificationRule) generatedapi.ListNotificationRulesResponse {
	apiRules := make([]generatedapi.NotificationRule, len(rules))
	for i, rule := range rules {
		apiRules[i] = DomainNotificationRuleToAPI(rule)
	}

	return generatedapi.ListNotificationRulesResponse{
		NotificationRules: apiRules,
	}
}

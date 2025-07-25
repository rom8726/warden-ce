package dto

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func strPtr(v string) *string { return &v }
func boolPtr(v bool) *bool    { return &v }

func TestDomainNotificationSettingToAPI(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.NotificationSetting
		expected generatedapi.NotificationSetting
	}{
		{
			name: "valid conversion",
			input: domain.NotificationSetting{
				ID:        1,
				ProjectID: 101,
				Type:      "email",
				Config:    json.RawMessage(`{"key":"value"}`),
				Enabled:   true,
				CreatedAt: time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC),
				UpdatedAt: time.Date(2023, 2, 3, 4, 5, 6, 0, time.UTC),
			},
			expected: generatedapi.NotificationSetting{
				ID:        1,
				ProjectID: 101,
				Type:      "email",
				Config:    `{"key":"value"}`,
				Enabled:   true,
				CreatedAt: time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC),
				UpdatedAt: time.Date(2023, 2, 3, 4, 5, 6, 0, time.UTC),
			},
		},
		{
			name: "empty config",
			input: domain.NotificationSetting{
				ID:        2,
				ProjectID: 0,
				Type:      "telegram",
				Config:    json.RawMessage(``),
				Enabled:   false,
				CreatedAt: time.Date(2023, 3, 2, 1, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 3, 3, 1, 0, 0, 0, time.UTC),
			},
			expected: generatedapi.NotificationSetting{
				ID:        2,
				ProjectID: 0,
				Type:      "telegram",
				Config:    ``,
				Enabled:   false,
				CreatedAt: time.Date(2023, 3, 2, 1, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 3, 3, 1, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := DomainNotificationSettingToAPI(tc.input)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("ожидалось %+v, получено %+v", tc.expected, got)
			}
		})
	}
}

/* ───────────────────────────── NotificationRule ────────────────────────────── */

func TestDomainNotificationRuleToAPI(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    domain.NotificationRule
		expected generatedapi.NotificationRule
	}{
		{
			name: "все опциональные заданы",
			input: domain.NotificationRule{
				ID:                  11,
				NotificationSetting: 5,
				EventLevel:          "error",
				Fingerprint:         strPtr("fp1"),
				IsNewError:          boolPtr(true),
				IsRegression:        boolPtr(true),
				CreatedAt:           now,
			},
			expected: generatedapi.NotificationRule{
				ID:                    11,
				NotificationSettingID: 5,
				EventLevel:            generatedapi.NewOptNilString("error"),
				Fingerprint: func() generatedapi.OptNilString {
					v := generatedapi.OptNilString{}
					v.Value, v.Set = "fp1", true
					return v
				}(),
				IsNewError: func() generatedapi.OptNilBool {
					v := generatedapi.OptNilBool{}
					v.Value, v.Set = true, true
					return v
				}(),
				IsRegression: generatedapi.NewOptNilBool(true),
				CreatedAt:    now,
			},
		},
		{
			name: "опциональные отсутствуют",
			input: domain.NotificationRule{
				ID:                  12,
				NotificationSetting: 6,
				EventLevel:          "warning",
				Fingerprint:         nil,
				IsNewError:          nil,
				IsRegression:        boolPtr(false),
				CreatedAt:           now,
			},
			expected: generatedapi.NotificationRule{
				ID:                    12,
				NotificationSettingID: 6,
				EventLevel:            generatedapi.NewOptNilString("warning"),
				// Fingerprint и IsNewError без Set
				Fingerprint:  generatedapi.OptNilString{},
				IsNewError:   generatedapi.OptNilBool{},
				IsRegression: generatedapi.NewOptNilBool(false),
				CreatedAt:    now,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := DomainNotificationRuleToAPI(tc.input)

			// Сравниваем поля по-отдельности, чтобы избежать ошибок из-за private-полей времени.
			if got.ID != tc.expected.ID ||
				got.NotificationSettingID != tc.expected.NotificationSettingID ||
				got.EventLevel.Value != tc.expected.EventLevel.Value ||
				got.EventLevel.Set != tc.expected.EventLevel.Set ||
				got.Fingerprint.Value != tc.expected.Fingerprint.Value ||
				got.Fingerprint.Set != tc.expected.Fingerprint.Set ||
				got.IsNewError.Value != tc.expected.IsNewError.Value ||
				got.IsNewError.Set != tc.expected.IsNewError.Set ||
				got.IsRegression.Value != tc.expected.IsRegression.Value ||
				got.IsRegression.Set != tc.expected.IsRegression.Set ||
				!got.CreatedAt.Equal(tc.expected.CreatedAt) {
				t.Fatalf("ожидалось %+v, получено %+v", tc.expected, got)
			}
		})
	}
}

/* ─────────────────────────── MakeNotificationRuleDTO ───────────────────────── */

func TestMakeNotificationRuleDTO(t *testing.T) {
	settingID := domain.NotificationSettingID(99)

	tests := []struct {
		name     string
		request  generatedapi.CreateNotificationRuleRequest
		expected domain.NotificationRuleDTO
	}{
		{
			name: "переданы все поля",
			request: generatedapi.CreateNotificationRuleRequest{
				EventLevel:   generatedapi.NewOptNilString("info"),
				Fingerprint:  generatedapi.NewOptNilString("xyz123"),
				IsNewError:   generatedapi.NewOptNilBool(false),
				IsRegression: generatedapi.NewOptNilBool(true),
			},
			expected: domain.NotificationRuleDTO{
				NotificationSetting: settingID,
				EventLevel:          "info",
				Fingerprint:         strPtr("xyz123"),
				IsNewError:          boolPtr(false),
				IsRegression:        boolPtr(true),
			},
		},
		{
			name: "передан только eventLevel",
			request: generatedapi.CreateNotificationRuleRequest{
				EventLevel: generatedapi.NewOptNilString("error"),
				// другие поля остаются unset
			},
			expected: domain.NotificationRuleDTO{
				NotificationSetting: settingID,
				EventLevel:          "error",
				Fingerprint:         nil,
				IsNewError:          nil,
				IsRegression:        nil,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := MakeNotificationRuleDTO(&tc.request, settingID)

			if got.NotificationSetting != tc.expected.NotificationSetting ||
				got.EventLevel != tc.expected.EventLevel ||
				(got.Fingerprint == nil) != (tc.expected.Fingerprint == nil) ||
				(got.IsNewError == nil) != (tc.expected.IsNewError == nil) ||
				(got.Fingerprint != nil && *got.Fingerprint != *tc.expected.Fingerprint) ||
				(got.IsNewError != nil && *got.IsNewError != *tc.expected.IsNewError) ||
				(got.IsRegression != nil && *got.IsRegression != *tc.expected.IsRegression) {
				t.Fatalf("ожидалось %+v, получено %+v", tc.expected, got)
			}
		})
	}
}

func TestMakeNotificationSettingDTO(t *testing.T) {
	tests := []struct {
		name      string
		req       generatedapi.CreateNotificationSettingRequest
		projectID domain.ProjectID
		expected  domain.NotificationSettingDTO
	}{
		{
			name: "Valid input",
			req: generatedapi.CreateNotificationSettingRequest{
				Type:    "email",
				Config:  `{"key":"value"}`,
				Enabled: generatedapi.NewOptBool(true),
			},
			projectID: 123,
			expected: domain.NotificationSettingDTO{
				ProjectID: 123,
				Type:      "email",
				Config:    json.RawMessage(`{"key":"value"}`),
				Enabled:   true,
			},
		},
		{
			name: "Missing required fields",
			req: generatedapi.CreateNotificationSettingRequest{
				Type:    "",
				Config:  ``,
				Enabled: generatedapi.NewOptBool(false),
			},
			projectID: 456,
			expected: domain.NotificationSettingDTO{
				ProjectID: 456,
				Type:      "",
				Config:    json.RawMessage(``),
				Enabled:   false,
			},
		},
		{
			name: "Empty Config and Type",
			req: generatedapi.CreateNotificationSettingRequest{
				Type:    "",
				Config:  ``,
				Enabled: generatedapi.NewOptBool(true),
			},
			projectID: 789,
			expected: domain.NotificationSettingDTO{
				ProjectID: 789,
				Type:      "",
				Config:    json.RawMessage(``),
				Enabled:   true,
			},
		},
		{
			name: "Nil Enabled value",
			req: generatedapi.CreateNotificationSettingRequest{
				Type:    "push",
				Config:  `{"anotherKey":"anotherValue"}`,
				Enabled: generatedapi.OptBool{},
			},
			projectID: 101,
			expected: domain.NotificationSettingDTO{
				ProjectID: 101,
				Type:      "push",
				Config:    json.RawMessage(`{"anotherKey":"anotherValue"}`),
				Enabled:   false,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := MakeNotificationSettingDTO(&tc.req, tc.projectID)
			if result.ProjectID != tc.expected.ProjectID ||
				result.Type != tc.expected.Type ||
				string(result.Config) != string(tc.expected.Config) ||
				result.Enabled != tc.expected.Enabled {
				t.Errorf("unexpected result: got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestUpdateNotificationSettingFromRequest(t *testing.T) {
	tests := []struct {
		name     string
		setting  domain.NotificationSetting
		req      generatedapi.UpdateNotificationSettingRequest
		expected domain.NotificationSetting
	}{
		{
			name: "Update all fields",
			setting: domain.NotificationSetting{
				ID:      1,
				Type:    "email",
				Enabled: true,
				Config:  json.RawMessage(`{"oldKey":"oldValue"}`),
			},
			req: generatedapi.UpdateNotificationSettingRequest{
				Type:    generatedapi.NewOptString("push"),
				Enabled: generatedapi.NewOptBool(false),
				Config:  generatedapi.NewOptString(`{"newKey":"newValue"}`),
			},
			expected: domain.NotificationSetting{
				ID:      1,
				Type:    "push",
				Enabled: false,
				Config:  json.RawMessage(`{"newKey":"newValue"}`),
			},
		},
		{
			name: "Update only Type",
			setting: domain.NotificationSetting{
				ID:      2,
				Type:    "email",
				Enabled: true,
				Config:  json.RawMessage(`{"key":"value"}`),
			},
			req: generatedapi.UpdateNotificationSettingRequest{
				Type: generatedapi.NewOptString("webhook"),
			},
			expected: domain.NotificationSetting{
				ID:      2,
				Type:    "webhook",
				Enabled: true,
				Config:  json.RawMessage(`{"key":"value"}`),
			},
		},
		{
			name: "Update only Enabled",
			setting: domain.NotificationSetting{
				ID:      3,
				Type:    "sms",
				Enabled: false,
				Config:  json.RawMessage(`{"key":"value"}`),
			},
			req: generatedapi.UpdateNotificationSettingRequest{
				Enabled: generatedapi.NewOptBool(true),
			},
			expected: domain.NotificationSetting{
				ID:      3,
				Type:    "sms",
				Enabled: true,
				Config:  json.RawMessage(`{"key":"value"}`),
			},
		},
		{
			name: "Unset Config",
			setting: domain.NotificationSetting{
				ID:      4,
				Type:    "email",
				Enabled: true,
				Config:  json.RawMessage(`{"key":"value"}`),
			},
			req: generatedapi.UpdateNotificationSettingRequest{
				Config: generatedapi.OptString{},
			},
			expected: domain.NotificationSetting{
				ID:      4,
				Type:    "email",
				Enabled: true,
				Config:  json.RawMessage(`{"key":"value"}`),
			},
		},
		{
			name: "Unset Type and Enabled",
			setting: domain.NotificationSetting{
				ID:      5,
				Type:    "push",
				Enabled: false,
				Config:  json.RawMessage(`{"key":"value"}`),
			},
			req: generatedapi.UpdateNotificationSettingRequest{
				Type:    generatedapi.OptString{},
				Enabled: generatedapi.OptBool{},
			},
			expected: domain.NotificationSetting{
				ID:      5,
				Type:    "push",
				Enabled: false,
				Config:  json.RawMessage(`{"key":"value"}`),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := UpdateNotificationSettingFromRequest(tc.setting, &tc.req)
			if result.ID != tc.expected.ID ||
				result.Type != tc.expected.Type ||
				string(result.Config) != string(tc.expected.Config) ||
				result.Enabled != tc.expected.Enabled {
				t.Errorf("unexpected result: got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestUpdateNotificationRuleFromRequest(t *testing.T) {
	t.SkipNow()

	tests := []struct {
		name     string
		rule     domain.NotificationRule
		req      generatedapi.UpdateNotificationRuleRequest
		expected domain.NotificationRule
	}{
		{
			name: "Update all fields",
			rule: domain.NotificationRule{
				ID:                  1,
				NotificationSetting: 123,
				EventLevel:          domain.IssueLevel("warning"),
				Fingerprint:         strPtr("oldFingerprint"),
				IsNewError:          boolPtr(false),
				IsRegression:        boolPtr(false),
			},
			req: generatedapi.UpdateNotificationRuleRequest{
				EventLevel:   generatedapi.NewOptNilString("error"),
				Fingerprint:  generatedapi.NewOptNilString("newFingerprint"),
				IsNewError:   generatedapi.NewOptNilBool(true),
				IsRegression: generatedapi.NewOptNilBool(true),
			},
			expected: domain.NotificationRule{
				ID:                  1,
				NotificationSetting: 123,
				EventLevel:          domain.IssueLevel("error"),
				Fingerprint:         strPtr("newFingerprint"),
				IsNewError:          boolPtr(true),
				IsRegression:        boolPtr(true),
			},
		},
		{
			name: "Update only EventLevel",
			rule: domain.NotificationRule{
				ID:                  2,
				NotificationSetting: 456,
				EventLevel:          domain.IssueLevel("info"),
				Fingerprint:         strPtr("sameFingerprint"),
				IsNewError:          boolPtr(false),
				IsRegression:        boolPtr(false),
			},
			req: generatedapi.UpdateNotificationRuleRequest{
				EventLevel: generatedapi.NewOptNilString("critical"),
			},
			expected: domain.NotificationRule{
				ID:                  2,
				NotificationSetting: 456,
				EventLevel:          domain.IssueLevel("critical"),
				Fingerprint:         strPtr("sameFingerprint"),
				IsNewError:          boolPtr(false),
				IsRegression:        boolPtr(false),
			},
		},
		{
			name: "Update only IsNewError",
			rule: domain.NotificationRule{
				ID:                  3,
				NotificationSetting: 789,
				EventLevel:          domain.IssueLevel("error"),
				Fingerprint:         strPtr("sameFingerprint"),
				IsNewError:          boolPtr(false),
				IsRegression:        boolPtr(true),
			},
			req: generatedapi.UpdateNotificationRuleRequest{
				IsNewError: generatedapi.NewOptNilBool(true),
			},
			expected: domain.NotificationRule{
				ID:                  3,
				NotificationSetting: 789,
				EventLevel:          domain.IssueLevel("error"),
				Fingerprint:         strPtr("sameFingerprint"),
				IsNewError:          boolPtr(true),
				IsRegression:        boolPtr(true),
			},
		},
		{
			name: "Unset Fingerprint and IsRegression",
			rule: domain.NotificationRule{
				ID:                  4,
				NotificationSetting: 101,
				EventLevel:          domain.IssueLevel("info"),
				Fingerprint:         strPtr("oldValue"),
				IsNewError:          boolPtr(true),
				IsRegression:        boolPtr(true),
			},
			req: generatedapi.UpdateNotificationRuleRequest{
				Fingerprint:  generatedapi.OptNilString{},
				IsRegression: generatedapi.OptNilBool{},
			},
			expected: domain.NotificationRule{
				ID:                  4,
				NotificationSetting: 101,
				EventLevel:          domain.IssueLevel("info"),
				Fingerprint:         strPtr("oldValue"),
				IsNewError:          boolPtr(true),
				IsRegression:        boolPtr(true),
			},
		},
		{
			name: "No updates (all nil values)",
			rule: domain.NotificationRule{
				ID:                  5,
				NotificationSetting: 202,
				EventLevel:          domain.IssueLevel("warning"),
				Fingerprint:         strPtr("unchanged"),
				IsNewError:          boolPtr(false),
				IsRegression:        boolPtr(false),
			},
			req: generatedapi.UpdateNotificationRuleRequest{
				EventLevel:   generatedapi.OptNilString{},
				Fingerprint:  generatedapi.OptNilString{},
				IsNewError:   generatedapi.OptNilBool{},
				IsRegression: generatedapi.OptNilBool{},
			},
			expected: domain.NotificationRule{
				ID:                  5,
				NotificationSetting: 202,
				EventLevel:          domain.IssueLevel("warning"),
				Fingerprint:         strPtr("unchanged"),
				IsNewError:          boolPtr(false),
				IsRegression:        boolPtr(false),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := UpdateNotificationRuleFromRequest(tc.rule, &tc.req)
			if result.ID != tc.expected.ID ||
				result.NotificationSetting != tc.expected.NotificationSetting ||
				result.EventLevel != tc.expected.EventLevel ||
				result.Fingerprint != tc.expected.Fingerprint ||
				result.IsNewError != tc.expected.IsNewError ||
				result.IsRegression != tc.expected.IsRegression {
				t.Errorf("unexpected result: got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestMakeListNotificationSettingsResponse(t *testing.T) {
	tests := []struct {
		name     string
		settings []domain.NotificationSetting
		expected generatedapi.ListNotificationSettingsResponse
	}{
		{
			name: "Valid input",
			settings: []domain.NotificationSetting{
				{
					ID:        1,
					ProjectID: 101,
					Type:      "email",
					Config:    json.RawMessage(`{"key":"value"}`),
					Enabled:   true,
					CreatedAt: time.Date(2025, 5, 13, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 5, 14, 15, 0, 0, 0, time.UTC),
				},
				{
					ID:        2,
					ProjectID: 102,
					Type:      "sms",
					Config:    nil,
					Enabled:   false,
					CreatedAt: time.Date(2025, 6, 1, 9, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 6, 2, 14, 0, 0, 0, time.UTC),
				},
			},
			expected: generatedapi.ListNotificationSettingsResponse{
				NotificationSettings: []generatedapi.NotificationSetting{
					{
						ID:        1,
						ProjectID: 101,
						Type:      "email",
						Config:    `{"key":"value"}`,
						Enabled:   true,
						CreatedAt: time.Date(2025, 5, 13, 10, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2025, 5, 14, 15, 0, 0, 0, time.UTC),
					},
					{
						ID:        2,
						ProjectID: 102,
						Type:      "sms",
						Config:    ``,
						Enabled:   false,
						CreatedAt: time.Date(2025, 6, 1, 9, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2025, 6, 2, 14, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			name:     "Empty input",
			settings: []domain.NotificationSetting{},
			expected: generatedapi.ListNotificationSettingsResponse{
				NotificationSettings: []generatedapi.NotificationSetting{},
			},
		},
		{
			name: "Mixed settings",
			settings: []domain.NotificationSetting{
				{
					ID:        3,
					ProjectID: 103,
					Type:      "push",
					Config:    json.RawMessage(`{"anotherKey":"anotherValue"}`),
					Enabled:   false,
					CreatedAt: time.Date(2025, 7, 1, 8, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 7, 2, 18, 0, 0, 0, time.UTC),
				},
			},
			expected: generatedapi.ListNotificationSettingsResponse{
				NotificationSettings: []generatedapi.NotificationSetting{
					{
						ID:        3,
						ProjectID: 103,
						Type:      "push",
						Config:    `{"anotherKey":"anotherValue"}`,
						Enabled:   false,
						CreatedAt: time.Date(2025, 7, 1, 8, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2025, 7, 2, 18, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := MakeListNotificationSettingsResponse(tc.settings)
			if len(result.NotificationSettings) != len(tc.expected.NotificationSettings) {
				t.Fatalf("unexpected length of notification settings: got %d, want %d", len(result.NotificationSettings), len(tc.expected.NotificationSettings))
			}
			for i, got := range result.NotificationSettings {
				want := tc.expected.NotificationSettings[i]
				if got.ID != want.ID ||
					got.ProjectID != want.ProjectID ||
					got.Type != want.Type ||
					got.Config != want.Config ||
					got.Enabled != want.Enabled ||
					!got.CreatedAt.Equal(want.CreatedAt) ||
					!got.UpdatedAt.Equal(want.UpdatedAt) {
					t.Errorf("unexpected result at index %d: got %+v, want %+v", i, got, want)
				}
			}
		})
	}
}

func TestMakeListNotificationRulesResponse(t *testing.T) {
	t.SkipNow()

	tests := []struct {
		name     string
		rules    []domain.NotificationRule
		expected generatedapi.ListNotificationRulesResponse
	}{
		{
			name: "Valid input",
			rules: []domain.NotificationRule{
				{
					ID:                  1,
					NotificationSetting: 101,
					EventLevel:          domain.IssueLevel("error"),
					Fingerprint:         strPtr("fingerprint123"),
					IsNewError:          boolPtr(true),
					IsRegression:        boolPtr(false),
					CreatedAt:           time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					ID:                  2,
					NotificationSetting: 102,
					EventLevel:          domain.IssueLevel("warning"),
					Fingerprint:         nil,
					IsNewError:          boolPtr(false),
					IsRegression:        boolPtr(true),
					CreatedAt:           time.Date(2025, 6, 2, 12, 0, 0, 0, time.UTC),
				},
			},
			expected: generatedapi.ListNotificationRulesResponse{
				NotificationRules: []generatedapi.NotificationRule{
					{
						ID:                    1,
						NotificationSettingID: 101,
						EventLevel:            generatedapi.NewOptNilString("error"),
						Fingerprint:           generatedapi.NewOptNilString("fingerprint123"),
						IsNewError:            generatedapi.NewOptNilBool(true),
						IsRegression:          generatedapi.NewOptNilBool(false),
						CreatedAt:             time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC),
					},
					{
						ID:                    2,
						NotificationSettingID: 102,
						EventLevel:            generatedapi.NewOptNilString("warning"),
						Fingerprint:           generatedapi.NewOptNilString(""),
						IsNewError:            generatedapi.NewOptNilBool(false),
						IsRegression:          generatedapi.NewOptNilBool(true),
						CreatedAt:             time.Date(2025, 6, 2, 12, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			name:  "Empty input",
			rules: []domain.NotificationRule{},
			expected: generatedapi.ListNotificationRulesResponse{
				NotificationRules: []generatedapi.NotificationRule{},
			},
		},
		{
			name: "Mixed input",
			rules: []domain.NotificationRule{
				{
					ID:                  3,
					NotificationSetting: 103,
					EventLevel:          domain.IssueLevel("info"),
					Fingerprint:         strPtr("mixedFingerprint"),
					IsNewError:          boolPtr(false),
					IsRegression:        boolPtr(true),
					CreatedAt:           time.Date(2025, 7, 1, 15, 0, 0, 0, time.UTC),
				},
			},
			expected: generatedapi.ListNotificationRulesResponse{
				NotificationRules: []generatedapi.NotificationRule{
					{
						ID:                    3,
						NotificationSettingID: 103,
						EventLevel:            generatedapi.NewOptNilString("info"),
						Fingerprint:           generatedapi.NewOptNilString("mixedFingerprint"),
						IsNewError:            generatedapi.NewOptNilBool(false),
						IsRegression:          generatedapi.NewOptNilBool(true),
						CreatedAt:             time.Date(2025, 7, 1, 15, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := MakeListNotificationRulesResponse(tc.rules)
			if len(result.NotificationRules) != len(tc.expected.NotificationRules) {
				t.Fatalf("unexpected length of notification rules: got %d, want %d", len(result.NotificationRules), len(tc.expected.NotificationRules))
			}
			for i, got := range result.NotificationRules {
				want := tc.expected.NotificationRules[i]
				if got.ID != want.ID ||
					got.NotificationSettingID != want.NotificationSettingID ||
					got.EventLevel != want.EventLevel ||
					got.Fingerprint != want.Fingerprint ||
					got.IsNewError != want.IsNewError ||
					got.IsRegression != want.IsRegression ||
					!got.CreatedAt.Equal(want.CreatedAt) {
					t.Errorf("unexpected result at index %d: got %+v, want %+v", i, got, want)
				}
			}
		})
	}
}

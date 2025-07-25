package notifications

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/issue-notificator/contract"
	mockdb "github.com/rom8726/warden/test_mocks/pkg/db"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestTakePendingNotificationsWithSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
			mockNotificationSettingsRepo *mockcontract.MockNotificationSettingsRepository,
			mockIssuesRepo *mockcontract.MockIssuesRepository,
		)
		limit          uint
		expectedResult []domain.NotificationWithSettings
		expectedError  bool
		errorContains  string
	}{
		{
			name: "Success - Multiple notifications and settings",
			setupMocks: func(
				mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
				mockNotificationSettingsRepo *mockcontract.MockNotificationSettingsRepository,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
			) {
				// Setup notifications
				notifications := []domain.Notification{
					{
						ID:        1,
						ProjectID: 100,
						IssueID:   1000,
						Level:     domain.IssueLevelError,
						IsNew:     true,
						Status:    domain.NotificationStatusPending,
						CreatedAt: time.Now(),
					},
					{
						ID:        2,
						ProjectID: 100,
						IssueID:   1001,
						Level:     domain.IssueLevelWarning,
						IsNew:     false,
						Status:    domain.NotificationStatusPending,
						CreatedAt: time.Now(),
					},
					{
						ID:        3,
						ProjectID: 200,
						IssueID:   2000,
						Level:     domain.IssueLevelError,
						IsNew:     true,
						Status:    domain.NotificationStatusPending,
						CreatedAt: time.Now(),
					},
				}

				mockNotificationsQueueRepo.EXPECT().TakePending(
					mock.Anything,
					uint(10),
				).Return(notifications, nil)

				// Setup settings for project 100
				settingsProject100 := []domain.NotificationSetting{
					{
						ID:        1,
						ProjectID: 100,
						Type:      domain.NotificationTypeEmail,
						Enabled:   true,
						CreatedAt: time.Now(),
						Rules: []domain.NotificationRule{
							{
								ID:                  1,
								NotificationSetting: 1,
								EventLevel:          domain.IssueLevelError,
								IsNewError:          boolPtr(true),
							},
						},
					},
					{
						ID:        2,
						ProjectID: 100,
						Type:      domain.NotificationTypeSlack,
						Enabled:   true,
						CreatedAt: time.Now(),
						Rules: []domain.NotificationRule{
							{
								ID:                  2,
								NotificationSetting: 2,
								EventLevel:          domain.IssueLevelWarning,
								IsNewError:          boolPtr(false),
							},
						},
					},
				}

				mockNotificationSettingsRepo.EXPECT().ListSettings(
					mock.Anything,
					domain.ProjectID(100),
				).Return(settingsProject100, nil)

				// Setup settings for project 200
				settingsProject200 := []domain.NotificationSetting{
					{
						ID:        3,
						ProjectID: 200,
						Type:      domain.NotificationTypeTelegram,
						Enabled:   true,
						CreatedAt: time.Now(),
						Rules: []domain.NotificationRule{
							{
								ID:                  3,
								NotificationSetting: 3,
								EventLevel:          domain.IssueLevelError,
								IsNewError:          boolPtr(true),
							},
						},
					},
				}

				mockNotificationSettingsRepo.EXPECT().ListSettings(
					mock.Anything,
					domain.ProjectID(200),
				).Return(settingsProject200, nil)
			},
			limit: 10,
			expectedResult: []domain.NotificationWithSettings{
				{
					Notification: domain.Notification{
						ID:        1,
						ProjectID: 100,
						IssueID:   1000,
						Level:     domain.IssueLevelError,
						IsNew:     true,
						Status:    domain.NotificationStatusPending,
					},
					Settings: []domain.NotificationSetting{
						{
							ID:        1,
							ProjectID: 100,
							Type:      domain.NotificationTypeEmail,
							Enabled:   true,
							Rules: []domain.NotificationRule{
								{
									ID:                  1,
									NotificationSetting: 1,
									EventLevel:          domain.IssueLevelError,
									IsNewError:          boolPtr(true),
								},
							},
						},
						{
							ID:        2,
							ProjectID: 100,
							Type:      domain.NotificationTypeSlack,
							Enabled:   true,
							Rules: []domain.NotificationRule{
								{
									ID:                  2,
									NotificationSetting: 2,
									EventLevel:          domain.IssueLevelWarning,
									IsNewError:          boolPtr(false),
								},
							},
						},
					},
				},
				{
					Notification: domain.Notification{
						ID:        2,
						ProjectID: 100,
						IssueID:   1001,
						Level:     domain.IssueLevelWarning,
						IsNew:     false,
						Status:    domain.NotificationStatusPending,
					},
					Settings: []domain.NotificationSetting{
						{
							ID:        1,
							ProjectID: 100,
							Type:      domain.NotificationTypeEmail,
							Enabled:   true,
							Rules: []domain.NotificationRule{
								{
									ID:                  1,
									NotificationSetting: 1,
									EventLevel:          domain.IssueLevelError,
									IsNewError:          boolPtr(true),
								},
							},
						},
						{
							ID:        2,
							ProjectID: 100,
							Type:      domain.NotificationTypeSlack,
							Enabled:   true,
							Rules: []domain.NotificationRule{
								{
									ID:                  2,
									NotificationSetting: 2,
									EventLevel:          domain.IssueLevelWarning,
									IsNewError:          boolPtr(false),
								},
							},
						},
					},
				},
				{
					Notification: domain.Notification{
						ID:        3,
						ProjectID: 200,
						IssueID:   2000,
						Level:     domain.IssueLevelError,
						IsNew:     true,
						Status:    domain.NotificationStatusPending,
					},
					Settings: []domain.NotificationSetting{
						{
							ID:        3,
							ProjectID: 200,
							Type:      domain.NotificationTypeTelegram,
							Enabled:   true,
							Rules: []domain.NotificationRule{
								{
									ID:                  3,
									NotificationSetting: 3,
									EventLevel:          domain.IssueLevelError,
									IsNewError:          boolPtr(true),
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Success - No notifications",
			setupMocks: func(
				mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
				mockNotificationSettingsRepo *mockcontract.MockNotificationSettingsRepository,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
			) {
				mockNotificationsQueueRepo.EXPECT().TakePending(
					mock.Anything,
					uint(10),
				).Return([]domain.Notification{}, nil)
			},
			limit:          10,
			expectedResult: []domain.NotificationWithSettings{},
			expectedError:  false,
		},
		{
			name: "Error - TakePending fails",
			setupMocks: func(
				mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
				mockNotificationSettingsRepo *mockcontract.MockNotificationSettingsRepository,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
			) {
				mockNotificationsQueueRepo.EXPECT().TakePending(
					mock.Anything,
					uint(10),
				).Return(nil, errors.New("database error"))
			},
			limit:          10,
			expectedResult: nil,
			expectedError:  true,
			errorContains:  "take pending notifications",
		},
		{
			name: "Error - ListSettings fails",
			setupMocks: func(
				mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
				mockNotificationSettingsRepo *mockcontract.MockNotificationSettingsRepository,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
			) {
				notifications := []domain.Notification{
					{
						ID:        1,
						ProjectID: 100,
						IssueID:   1000,
						Level:     domain.IssueLevelError,
						IsNew:     true,
						Status:    domain.NotificationStatusPending,
						CreatedAt: time.Now(),
					},
				}

				mockNotificationsQueueRepo.EXPECT().TakePending(
					mock.Anything,
					uint(10),
				).Return(notifications, nil)

				mockNotificationSettingsRepo.EXPECT().ListSettings(
					mock.Anything,
					domain.ProjectID(100),
				).Return(nil, errors.New("database error"))
			},
			limit:          10,
			expectedResult: nil,
			expectedError:  true,
			errorContains:  "list notification settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockNotificationSettingsRepo := mockcontract.NewMockNotificationSettingsRepository(t)
			mockNotificationsQueueRepo := mockcontract.NewMockNotificationsQueueRepository(t)
			issuesRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(mockNotificationsQueueRepo, mockNotificationSettingsRepo, issuesRepo)

			// Create service
			service := New(
				mockTxManager,
				mockNotificationSettingsRepo,
				mockNotificationsQueueRepo,
				issuesRepo,
			)

			// Call method
			result, err := service.TakePendingNotificationsWithSettings(context.Background(), tt.limit)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.expectedResult), len(result))

				// Compare each notification with settings
				for i, expectedNotification := range tt.expectedResult {
					require.Equal(t, expectedNotification.ID, result[i].ID)
					require.Equal(t, expectedNotification.ProjectID, result[i].ProjectID)
					require.Equal(t, expectedNotification.IssueID, result[i].IssueID)
					require.Equal(t, expectedNotification.Level, result[i].Level)
					require.Equal(t, expectedNotification.IsNew, result[i].IsNew)
					require.Equal(t, expectedNotification.Status, result[i].Status)

					// Compare settings
					require.Equal(t, len(expectedNotification.Settings), len(result[i].Settings))
					for j, expectedSetting := range expectedNotification.Settings {
						require.Equal(t, expectedSetting.ID, result[i].Settings[j].ID)
						require.Equal(t, expectedSetting.ProjectID, result[i].Settings[j].ProjectID)
						require.Equal(t, expectedSetting.Type, result[i].Settings[j].Type)
						require.Equal(t, expectedSetting.Enabled, result[i].Settings[j].Enabled)

						// Compare rules
						require.Equal(t, len(expectedSetting.Rules), len(result[i].Settings[j].Rules))
						for k, expectedRule := range expectedSetting.Rules {
							require.Equal(t, expectedRule.ID, result[i].Settings[j].Rules[k].ID)
							require.Equal(t, expectedRule.NotificationSetting, result[i].Settings[j].Rules[k].NotificationSetting)
							require.Equal(t, expectedRule.EventLevel, result[i].Settings[j].Rules[k].EventLevel)
							require.Equal(t, expectedRule.IsNewError, result[i].Settings[j].Rules[k].IsNewError)
						}
					}
				}
			}
		})
	}
}

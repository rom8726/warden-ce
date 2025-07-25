package notificator

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/issue-notificator/contract"
	mocknotificator "github.com/rom8726/warden/test_mocks/internal_/issue-notificator/notificator"
	mockdb "github.com/rom8726/warden/test_mocks/pkg/db"
)

func newMockChannel(nt domain.NotificationType) *mocknotificator.MockChannel {
	ch := &mocknotificator.MockChannel{}
	ch.EXPECT().Type().Return(nt)
	return ch
}

func boolPtr(b bool) *bool {
	return &b
}

func TestProcessOutbox(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockNotificationsUseCase *mockcontract.MockNotificationsUseCase,
			mockIssuesRepo *mockcontract.MockIssuesRepository,
			mockProjectsRepo *mockcontract.MockProjectsRepository,
			mockEmailChannel *mocknotificator.MockChannel,
		)
		expectedSentCount int
		checkLogOutput    bool
	}{
		{
			name: "Success - Multiple notifications processed in batches",
			setupMocks: func(
				mockNotificationsUseCase *mockcontract.MockNotificationsUseCase,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEmailChannel *mocknotificator.MockChannel,
			) {
				// First batch
				notifications1 := []domain.NotificationWithSettings{
					{
						Notification: domain.Notification{
							ID:             domain.NotificationID(1),
							ProjectID:      domain.ProjectID(100),
							IssueID:        domain.IssueID(1000),
							Level:          domain.IssueLevelError,
							IsNew:          true,
							WasReactivated: false,
							Status:         domain.NotificationStatusPending,
						},
						Settings: []domain.NotificationSetting{
							{
								ID:        domain.NotificationSettingID(1),
								ProjectID: domain.ProjectID(100),
								Type:      domain.NotificationTypeEmail,
								Config:    json.RawMessage(`{"to": "user@example.com"}`),
								Enabled:   true,
								Rules: []domain.NotificationRule{
									{
										ID:                  domain.NotificationRuleID(1),
										NotificationSetting: domain.NotificationSettingID(1),
										EventLevel:          domain.IssueLevelError,
										IsNewError:          boolPtr(true),
										IsRegression:        boolPtr(false),
									},
								},
							},
						},
					},
				}

				// Second batch
				notifications2 := []domain.NotificationWithSettings{
					{
						Notification: domain.Notification{
							ID:             domain.NotificationID(2),
							ProjectID:      domain.ProjectID(100),
							IssueID:        domain.IssueID(1001),
							Level:          domain.IssueLevelError,
							IsNew:          true,
							WasReactivated: false,
							Status:         domain.NotificationStatusPending,
						},
						Settings: []domain.NotificationSetting{
							{
								ID:        domain.NotificationSettingID(2),
								ProjectID: domain.ProjectID(100),
								Type:      domain.NotificationTypeEmail,
								Config:    json.RawMessage(`{"to": "user@example.com"}`),
								Enabled:   true,
								Rules: []domain.NotificationRule{
									{
										ID:                  domain.NotificationRuleID(2),
										NotificationSetting: domain.NotificationSettingID(2),
										EventLevel:          domain.IssueLevelError,
										IsNewError:          boolPtr(true),
										IsRegression:        boolPtr(false),
									},
								},
							},
						},
					},
				}

				// Empty batch to end the loop
				notifications3 := []domain.NotificationWithSettings{}

				mockNotificationsUseCase.
					EXPECT().
					TakePendingNotificationsWithSettings(
						mock.Anything,
						uint(10),
					).
					Return(notifications1, nil).
					Once()

				mockNotificationsUseCase.
					EXPECT().
					TakePendingNotificationsWithSettings(
						mock.Anything,
						uint(10),
					).
					Return(notifications2, nil).
					Once()

				mockNotificationsUseCase.
					EXPECT().
					TakePendingNotificationsWithSettings(
						mock.Anything,
						uint(10),
					).
					Return(notifications3, nil).
					Once()

				// Setup expectations for first notification
				mockIssuesRepo.
					EXPECT().
					GetByID(mock.Anything, domain.IssueID(1000)).
					Return(domain.Issue{
						ID:        domain.IssueID(1000),
						ProjectID: domain.ProjectID(100),
					}, nil)

				mockProjectsRepo.
					EXPECT().
					GetByID(mock.Anything, domain.ProjectID(100)).
					Return(domain.Project{
						ID: domain.ProjectID(100),
					}, nil)

				mockEmailChannel.
					EXPECT().
					Send(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				mockNotificationsUseCase.
					EXPECT().
					MarkNotificationAsSent(mock.Anything, domain.NotificationID(1)).
					Return(nil)

				// Setup expectations for second notification
				mockIssuesRepo.
					EXPECT().
					GetByID(mock.Anything, domain.IssueID(1001)).
					Return(domain.Issue{
						ID:        domain.IssueID(1001),
						ProjectID: domain.ProjectID(100),
					}, nil)

				mockProjectsRepo.
					EXPECT().
					GetByID(mock.Anything, domain.ProjectID(100)).
					Return(domain.Project{
						ID: domain.ProjectID(100),
					}, nil)

				mockNotificationsUseCase.
					EXPECT().
					MarkNotificationAsSent(mock.Anything, domain.NotificationID(2)).
					Return(nil)
			},
			expectedSentCount: 2,
			checkLogOutput:    false,
		},
		{
			name: "Success - Notifications skipped when no settings match",
			setupMocks: func(
				mockNotificationsUseCase *mockcontract.MockNotificationsUseCase,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEmailChannel *mocknotificator.MockChannel,
			) {
				notifications := []domain.NotificationWithSettings{
					{
						Notification: domain.Notification{
							ID:             domain.NotificationID(1),
							ProjectID:      domain.ProjectID(100),
							IssueID:        domain.IssueID(1000),
							Level:          domain.IssueLevelError,
							IsNew:          true,
							WasReactivated: false,
							Status:         domain.NotificationStatusPending,
						},
						Settings: []domain.NotificationSetting{
							{
								ID:        domain.NotificationSettingID(1),
								ProjectID: domain.ProjectID(100),
								Type:      domain.NotificationTypeEmail,
								Config:    json.RawMessage(`{"to": "user@example.com"}`),
								Enabled:   false, // Disabled setting
								Rules: []domain.NotificationRule{
									{
										ID:                  domain.NotificationRuleID(1),
										NotificationSetting: domain.NotificationSettingID(1),
										EventLevel:          domain.IssueLevelError,
										IsNewError:          boolPtr(true),
										IsRegression:        boolPtr(false),
									},
								},
							},
						},
					},
				}

				// Empty batch to end the loop
				emptyNotifications := []domain.NotificationWithSettings{}

				mockNotificationsUseCase.
					EXPECT().
					TakePendingNotificationsWithSettings(
						mock.Anything,
						uint(10),
					).
					Return(notifications, nil).
					Once()

				mockNotificationsUseCase.
					EXPECT().
					TakePendingNotificationsWithSettings(
						mock.Anything,
						uint(10),
					).
					Return(emptyNotifications, nil).
					Once()

				mockIssuesRepo.
					EXPECT().
					GetByID(mock.Anything, domain.IssueID(1000)).
					Return(domain.Issue{
						ID:        domain.IssueID(1000),
						ProjectID: domain.ProjectID(100),
					}, nil)

				mockProjectsRepo.
					EXPECT().
					GetByID(mock.Anything, domain.ProjectID(100)).
					Return(domain.Project{
						ID: domain.ProjectID(100),
					}, nil)

				// Should be marked as skipped
				mockNotificationsUseCase.
					EXPECT().
					MarkNotificationAsSkipped(mock.Anything, domain.NotificationID(1), "no settings").
					Return(nil)
			},
			expectedSentCount: 0,
			checkLogOutput:    false,
		},
		{
			name: "Success - Notifications skipped when issue not found",
			setupMocks: func(
				mockNotificationsUseCase *mockcontract.MockNotificationsUseCase,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEmailChannel *mocknotificator.MockChannel,
			) {
				notifications := []domain.NotificationWithSettings{
					{
						Notification: domain.Notification{
							ID:             domain.NotificationID(1),
							ProjectID:      domain.ProjectID(100),
							IssueID:        domain.IssueID(1000),
							Level:          domain.IssueLevelError,
							IsNew:          true,
							WasReactivated: false,
							Status:         domain.NotificationStatusPending,
						},
						Settings: []domain.NotificationSetting{
							{
								ID:        domain.NotificationSettingID(1),
								ProjectID: domain.ProjectID(100),
								Type:      domain.NotificationTypeEmail,
								Config:    json.RawMessage(`{"to": "user@example.com"}`),
								Enabled:   true,
								Rules: []domain.NotificationRule{
									{
										ID:                  domain.NotificationRuleID(1),
										NotificationSetting: domain.NotificationSettingID(1),
										EventLevel:          domain.IssueLevelError,
										IsNewError:          boolPtr(true),
										IsRegression:        boolPtr(false),
									},
								},
							},
						},
					},
				}

				// Empty batch to end the loop
				emptyNotifications := []domain.NotificationWithSettings{}

				mockNotificationsUseCase.
					EXPECT().
					TakePendingNotificationsWithSettings(
						mock.Anything,
						uint(10),
					).
					Return(notifications, nil).
					Once()

				mockNotificationsUseCase.
					EXPECT().
					TakePendingNotificationsWithSettings(
						mock.Anything,
						uint(10),
					).
					Return(emptyNotifications, nil).
					Once()

				mockIssuesRepo.
					EXPECT().
					GetByID(mock.Anything, domain.IssueID(1000)).
					Return(domain.Issue{}, domain.ErrEntityNotFound)

				// Should be marked as skipped
				mockNotificationsUseCase.
					EXPECT().
					MarkNotificationAsSkipped(mock.Anything, domain.NotificationID(1), "issue not found").
					Return(nil)
			},
			expectedSentCount: 0,
			checkLogOutput:    false,
		},
		{
			name: "Success - Notifications skipped when project not found",
			setupMocks: func(
				mockNotificationsUseCase *mockcontract.MockNotificationsUseCase,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEmailChannel *mocknotificator.MockChannel,
			) {
				notifications := []domain.NotificationWithSettings{
					{
						Notification: domain.Notification{
							ID:             domain.NotificationID(1),
							ProjectID:      domain.ProjectID(100),
							IssueID:        domain.IssueID(1000),
							Level:          domain.IssueLevelError,
							IsNew:          true,
							WasReactivated: false,
							Status:         domain.NotificationStatusPending,
						},
						Settings: []domain.NotificationSetting{
							{
								ID:        domain.NotificationSettingID(1),
								ProjectID: domain.ProjectID(100),
								Type:      domain.NotificationTypeEmail,
								Config:    json.RawMessage(`{"to": "user@example.com"}`),
								Enabled:   true,
								Rules: []domain.NotificationRule{
									{
										ID:                  domain.NotificationRuleID(1),
										NotificationSetting: domain.NotificationSettingID(1),
										EventLevel:          domain.IssueLevelError,
										IsNewError:          boolPtr(true),
										IsRegression:        boolPtr(false),
									},
								},
							},
						},
					},
				}

				// Empty batch to end the loop
				emptyNotifications := []domain.NotificationWithSettings{}

				mockNotificationsUseCase.
					EXPECT().
					TakePendingNotificationsWithSettings(
						mock.Anything,
						uint(10),
					).
					Return(notifications, nil).
					Once()

				mockNotificationsUseCase.
					EXPECT().
					TakePendingNotificationsWithSettings(
						mock.Anything,
						uint(10),
					).
					Return(emptyNotifications, nil).
					Once()

				mockIssuesRepo.
					EXPECT().
					GetByID(mock.Anything, domain.IssueID(1000)).
					Return(domain.Issue{
						ID:        domain.IssueID(1000),
						ProjectID: domain.ProjectID(100),
					}, nil)

				mockProjectsRepo.
					EXPECT().
					GetByID(mock.Anything, domain.ProjectID(100)).
					Return(domain.Project{}, domain.ErrEntityNotFound)

				// Should be marked as skipped
				mockNotificationsUseCase.
					EXPECT().
					MarkNotificationAsSkipped(mock.Anything, domain.NotificationID(1), "project not found").
					Return(nil)
			},
			expectedSentCount: 0,
			checkLogOutput:    false,
		},
		{
			name: "Error - TakePendingNotificationsWithSettings fails",
			setupMocks: func(
				mockNotificationsUseCase *mockcontract.MockNotificationsUseCase,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEmailChannel *mocknotificator.MockChannel,
			) {
				mockNotificationsUseCase.
					EXPECT().
					TakePendingNotificationsWithSettings(
						mock.Anything,
						uint(10),
					).
					Return(nil, assert.AnError).
					Once()
			},
			expectedSentCount: 0,
			checkLogOutput:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockTxManager := mockdb.NewMockTxManager(t)
			mockNotificationsUseCase := &mockcontract.MockNotificationsUseCase{}
			mockIssuesRepo := &mockcontract.MockIssuesRepository{}
			mockProjectsRepo := &mockcontract.MockProjectsRepository{}
			mockEmailChannel := newMockChannel(domain.NotificationTypeEmail)

			tt.setupMocks(
				mockNotificationsUseCase,
				mockIssuesRepo,
				mockProjectsRepo,
				mockEmailChannel,
			)

			svc := New(
				[]Channel{mockEmailChannel},
				mockTxManager,
				mockNotificationsUseCase,
				mockIssuesRepo,
				mockProjectsRepo,
				4, // workerCount
			)

			svc.batchSize = 10
			svc.ProcessOutbox(context.Background())

			mockNotificationsUseCase.AssertExpectations(t)
			mockIssuesRepo.AssertExpectations(t)
			mockProjectsRepo.AssertExpectations(t)
			mockEmailChannel.AssertExpectations(t)
		})
	}
}

func TestProcessOutboxParallel(t *testing.T) {
	t.Parallel()

	// Create multiple notifications to test parallel processing
	notifications := make([]domain.NotificationWithSettings, 10)
	for i := 0; i < 10; i++ {
		notifications[i] = domain.NotificationWithSettings{
			Notification: domain.Notification{
				ID:             domain.NotificationID(i + 1),
				ProjectID:      domain.ProjectID(100),
				IssueID:        domain.IssueID(1000 + i),
				Level:          domain.IssueLevelError,
				IsNew:          true,
				WasReactivated: false,
				Status:         domain.NotificationStatusPending,
			},
			Settings: []domain.NotificationSetting{
				{
					ID:        domain.NotificationSettingID(i + 1),
					ProjectID: domain.ProjectID(100),
					Type:      domain.NotificationTypeEmail,
					Config:    json.RawMessage(`{"to": "user@example.com"}`),
					Enabled:   true,
					Rules: []domain.NotificationRule{
						{
							ID:                  domain.NotificationRuleID(i + 1),
							NotificationSetting: domain.NotificationSettingID(i + 1),
							EventLevel:          domain.IssueLevelError,
							IsNewError:          boolPtr(true),
							IsRegression:        boolPtr(false),
						},
					},
				},
			},
		}
	}

	// Empty batch to end the loop
	emptyNotifications := []domain.NotificationWithSettings{}

	workerCounts := []int{1, 2, 4, 8}
	for _, workerCount := range workerCounts {
		t.Run(fmt.Sprintf("WorkerCount_%d", workerCount), func(t *testing.T) {
			t.Parallel()

			mockTxManager := mockdb.NewMockTxManager(t)
			mockNotificationsUseCase := &mockcontract.MockNotificationsUseCase{}
			mockIssuesRepo := &mockcontract.MockIssuesRepository{}
			mockProjectsRepo := &mockcontract.MockProjectsRepository{}
			mockEmailChannel := newMockChannel(domain.NotificationTypeEmail)

			mockNotificationsUseCase.
				EXPECT().
				TakePendingNotificationsWithSettings(
					mock.Anything,
					uint(10),
				).
				Return(notifications, nil).
				Once()

			mockNotificationsUseCase.
				EXPECT().
				TakePendingNotificationsWithSettings(
					mock.Anything,
					uint(10),
				).
				Return(emptyNotifications, nil).
				Once()

			for i := 0; i < 10; i++ {
				mockIssuesRepo.
					EXPECT().
					GetByID(mock.Anything, domain.IssueID(1000+i)).
					Return(domain.Issue{
						ID:        domain.IssueID(1000 + i),
						ProjectID: domain.ProjectID(100),
					}, nil)

				mockProjectsRepo.
					EXPECT().
					GetByID(mock.Anything, domain.ProjectID(100)).
					Return(domain.Project{
						ID: domain.ProjectID(100),
					}, nil)

				mockEmailChannel.
					EXPECT().
					Send(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				mockNotificationsUseCase.
					EXPECT().
					MarkNotificationAsSent(mock.Anything, domain.NotificationID(i+1)).
					Return(nil)
			}

			svc := New(
				[]Channel{mockEmailChannel},
				mockTxManager,
				mockNotificationsUseCase,
				mockIssuesRepo,
				mockProjectsRepo,
				workerCount,
			)

			svc.batchSize = 10
			svc.ProcessOutbox(context.Background())

			mockNotificationsUseCase.AssertExpectations(t)
			mockIssuesRepo.AssertExpectations(t)
			mockProjectsRepo.AssertExpectations(t)
			mockEmailChannel.AssertExpectations(t)
		})
	}
}

package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/envelope-consumer/contract"
	mockdb "github.com/rom8726/warden/test_mocks/pkg/db"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockTxManager := mockdb.NewMockTxManager(t)
	mockIssueRepo := mockcontract.NewMockIssuesRepository(t)
	mockEventRepo := mockcontract.NewMockEventRepository(t)
	mockNotificationsQueueRepo := mockcontract.NewMockNotificationsQueueRepository(t)
	mockReleaseRepo := mockcontract.NewMockReleaseRepository(t)
	mockIssueReleaseRepo := mockcontract.NewMockIssueReleasesRepository(t)
	cacheService := mockcontract.NewMockCacheService(t)

	// Create service
	service := New(
		mockTxManager,
		mockIssueRepo,
		mockEventRepo,
		mockReleaseRepo,
		mockNotificationsQueueRepo,
		mockIssueReleaseRepo,
		cacheService,
	)
	// Verify service was created correctly
	require.NotNil(t, service)
	require.Equal(t, mockTxManager, service.txManager)
	require.Equal(t, mockIssueRepo, service.issueRepo)
	require.Equal(t, mockEventRepo, service.eventRepo)
	require.Equal(t, mockNotificationsQueueRepo, service.notificationsQueueRepo)
}

func TestProcessEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		eventData  map[string]any
		setupMocks func(
			mockTxManager *mockdb.MockTxManager,
			mockIssueRepo *mockcontract.MockIssuesRepository,
			mockEventRepo *mockcontract.MockEventRepository,
			mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
			mockReleaseRepo *mockcontract.MockReleaseRepository,
			mockIssueReleaseRepo *mockcontract.MockIssueReleasesRepository,
			mockCacheService *mockcontract.MockCacheService,
		)
		expectedEventID     domain.EventID
		expectedError       bool
		expectedErrorString string
	}{
		{
			name: "Missing event_id",
			eventData: map[string]any{
				"message": "Test message",
			},
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssueRepo *mockcontract.MockIssuesRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockIssueReleaseRepo *mockcontract.MockIssueReleasesRepository,
				mockCacheService *mockcontract.MockCacheService,
			) {
				// No mocks needed for this case
			},
			expectedEventID:     "",
			expectedError:       true,
			expectedErrorString: "event_id is required",
		},
		{
			name: "Basic event processing",
			eventData: map[string]any{
				"event_id":  "test-event-id",
				"message":   "Test message",
				"level":     "error",
				"platform":  "python",
				"timestamp": time.Now().Format(time.RFC3339),
				"tags": map[string]any{
					"tag1": "value1",
					"tag2": "value2",
				},
				"environment": "production",
				"server_name": "test-server",
			},
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssueRepo *mockcontract.MockIssuesRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockIssueReleaseRepo *mockcontract.MockIssueReleasesRepository,
				mockCacheService *mockcontract.MockCacheService,
			) {
				// Setup transaction
				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Run(func(ctx context.Context, fn func(context.Context) error) {
						err := fn(ctx)
						require.NoError(t, err)
					}).Return(nil)

				// Setup event repo
				mockEventRepo.EXPECT().StoreWithFingerprints(mock.Anything, mock.AnythingOfType("*domain.Event")).
					Return(nil)

				// Setup issue repo
				mockIssueRepo.EXPECT().UpsertIssue(mock.Anything, mock.AnythingOfType("domain.IssueDTO")).
					Return(domain.IssueUpsertResult{
						ID:    123,
						IsNew: false,
					}, nil)

				// Setup release repo
				//mockReleaseRepo.On("Create", mock.Anything, mock.AnythingOfType("domain.ReleaseDTO")).
				//	Return(domain.ReleaseID(1), nil)

				// Setup issue release repo
				//mockIssueReleaseRepo.On("Create", mock.Anything, domain.IssueID(123), domain.ReleaseID(1), false).
				//	Return(nil)

				mockCacheService.EXPECT().GetOrCreateRelease(mock.Anything, domain.ProjectID(1), "unknown", mock.Anything).
					Return(domain.ReleaseID(1), nil)

				mockCacheService.EXPECT().GetOrCreateIssueRelease(mock.Anything, domain.IssueID(123), domain.ReleaseID(1), false, mock.Anything).
					Return(nil)
			},
			expectedEventID: "test-event-id",
			expectedError:   false,
		},
		{
			name: "Event with exception",
			eventData: map[string]any{
				"event_id": "test-event-id",
				"message":  "Test message",
				"level":    "error",
				"platform": "python",
				"exception": map[string]any{
					"type":  "ValueError",
					"value": "Invalid value",
					"stacktrace": map[string]any{
						"frames": []any{
							map[string]any{
								"filename": "test.py",
								"lineno":   10,
							},
						},
					},
				},
			},
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssueRepo *mockcontract.MockIssuesRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockIssueReleaseRepo *mockcontract.MockIssueReleasesRepository,
				mockCacheService *mockcontract.MockCacheService,
			) {
				// Setup transaction
				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Run(func(ctx context.Context, fn func(context.Context) error) {
						err := fn(ctx)
						require.NoError(t, err)
					}).Return(nil)

				// Setup event repo
				mockEventRepo.EXPECT().StoreWithFingerprints(mock.Anything, mock.AnythingOfType("*domain.Event")).
					Return(nil)

				// Setup issue repo
				mockIssueRepo.EXPECT().UpsertIssue(mock.Anything, mock.AnythingOfType("domain.IssueDTO")).
					Return(domain.IssueUpsertResult{
						ID:    123,
						IsNew: false,
					}, nil)

				// Setup release repo
				//mockReleaseRepo.On("Create", mock.Anything, mock.AnythingOfType("domain.ReleaseDTO")).
				//	Return(domain.ReleaseID(1), nil)

				// Setup issue release repo
				//mockIssueReleaseRepo.On("Create", mock.Anything, domain.IssueID(123), domain.ReleaseID(1), false).
				//	Return(nil)
				mockCacheService.EXPECT().GetOrCreateRelease(mock.Anything, domain.ProjectID(1), "unknown", mock.Anything).
					Return(domain.ReleaseID(1), nil)

				mockCacheService.EXPECT().GetOrCreateIssueRelease(mock.Anything, domain.IssueID(123), domain.ReleaseID(1), false, mock.Anything).
					Return(nil)
			},
			expectedEventID: "test-event-id",
			expectedError:   false,
		},
		{
			name: "New issue triggers notification",
			eventData: map[string]any{
				"event_id": "test-event-id",
				"message":  "Test message",
				"level":    "error",
			},
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssueRepo *mockcontract.MockIssuesRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockIssueReleaseRepo *mockcontract.MockIssueReleasesRepository,
				mockCacheService *mockcontract.MockCacheService,
			) {
				// Setup transaction
				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Run(func(ctx context.Context, fn func(context.Context) error) {
						err := fn(ctx)
						require.NoError(t, err)
					}).Return(nil)

				// Setup event repo
				mockEventRepo.EXPECT().StoreWithFingerprints(mock.Anything, mock.AnythingOfType("*domain.Event")).
					Return(nil)

				// Setup issue repo - return IsNew=true to trigger notification
				mockIssueRepo.EXPECT().UpsertIssue(mock.Anything, mock.AnythingOfType("domain.IssueDTO")).
					Return(domain.IssueUpsertResult{
						ID:    123,
						IsNew: true,
					}, nil)

				// Setup release repo
				//mockReleaseRepo.On("Create", mock.Anything, mock.AnythingOfType("domain.ReleaseDTO")).
				//	Return(domain.ReleaseID(1), nil)

				// Setup issue release repo
				//mockIssueReleaseRepo.On("Create", mock.Anything, domain.IssueID(123), domain.ReleaseID(1), true).
				//	Return(nil)

				mockCacheService.EXPECT().GetOrCreateRelease(mock.Anything, domain.ProjectID(1), "unknown", mock.Anything).
					Return(domain.ReleaseID(1), nil)

				mockCacheService.EXPECT().GetOrCreateIssueRelease(mock.Anything, domain.IssueID(123), domain.ReleaseID(1), true, mock.Anything).
					Return(nil)

				// Setup notifications queue repo
				mockNotificationsQueueRepo.EXPECT().AddNotification(
					mock.Anything,
					domain.ProjectID(1),
					domain.IssueID(123),
					domain.IssueLevel("error"),
					true,
					false).
					Return(nil)
			},
			expectedEventID: "test-event-id",
			expectedError:   false,
		},
		{
			name: "Error storing event",
			eventData: map[string]any{
				"event_id": "test-event-id",
				"message":  "Test message",
			},
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssueRepo *mockcontract.MockIssuesRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockNotificationsQueueRepo *mockcontract.MockNotificationsQueueRepository,
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockIssueReleaseRepo *mockcontract.MockIssueReleasesRepository,
				mockCacheService *mockcontract.MockCacheService,
			) {
				// Setup transaction to return error
				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Run(func(ctx context.Context, fn func(context.Context) error) {
						_ = fn(ctx)
					}).Return(errors.New("store event and issue: store event: database error"))

				// Setup event repo to return error
				mockEventRepo.EXPECT().StoreWithFingerprints(mock.Anything, mock.AnythingOfType("*domain.Event")).
					Return(errors.New("database error"))
			},
			expectedEventID:     "",
			expectedError:       true,
			expectedErrorString: "store event and issue: store event and issue: store event: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockIssueRepo := mockcontract.NewMockIssuesRepository(t)
			mockEventRepo := mockcontract.NewMockEventRepository(t)
			mockNotificationsQueueRepo := mockcontract.NewMockNotificationsQueueRepository(t)
			mockReleaseRepo := mockcontract.NewMockReleaseRepository(t)
			mockIssueReleaseRepo := mockcontract.NewMockIssueReleasesRepository(t)
			cacheService := mockcontract.NewMockCacheService(t)

			// Setup mocks
			tt.setupMocks(
				mockTxManager,
				mockIssueRepo,
				mockEventRepo,
				mockNotificationsQueueRepo,
				mockReleaseRepo,
				mockIssueReleaseRepo,
				cacheService,
			)

			// Create service
			service := New(
				mockTxManager,
				mockIssueRepo,
				mockEventRepo,
				mockReleaseRepo,
				mockNotificationsQueueRepo,
				mockIssueReleaseRepo,
				cacheService,
			)

			// Call the method
			eventID, err := service.ProcessEvent(context.Background(), 1, tt.eventData)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
				if tt.expectedErrorString != "" {
					require.Contains(t, err.Error(), tt.expectedErrorString)
				}
				require.Equal(t, tt.expectedEventID, eventID)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedEventID, eventID)
			}
		})
	}
}

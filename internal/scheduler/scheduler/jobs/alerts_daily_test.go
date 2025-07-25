package jobs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/scheduler/contract"
)

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestAlertsDailyJob_Run(t *testing.T) {
	today := time.Now().Truncate(24 * time.Hour)
	firstSeen := today.Add(-3 * 24 * time.Hour)
	regressAt := firstSeen.Add(24 * time.Hour)

	type testCase struct {
		name           string
		issues         []domain.IssueExtended
		listErr        error
		addNotifErr    error
		expectedQueued int
	}

	cases := []testCase{
		{
			name: "one unresolved, no notifications yet",
			issues: []domain.IssueExtended{
				{
					Issue: domain.Issue{
						ID:        1,
						ProjectID: 1,
						Level:     domain.IssueLevelError,
						CreatedAt: firstSeen,
						Status:    domain.IssueStatusUnresolved,
					},
				},
			},
			expectedQueued: 1,
		},
		{
			name:           "no unresolved issues",
			issues:         nil,
			expectedQueued: 0,
		},
		{
			name:           "repo returns error",
			listErr:        errors.New("fail"),
			expectedQueued: 0,
		},
		{
			name: "add notification fails",
			issues: []domain.IssueExtended{
				{
					Issue: domain.Issue{
						ID:        2,
						ProjectID: 2,
						Level:     domain.IssueLevelError,
						CreatedAt: firstSeen,
						Status:    domain.IssueStatusUnresolved,
					},
				},
			},
			addNotifErr:    errors.New("fail add"),
			expectedQueued: 0,
		},
		{
			name: "regress issue, no last_notification_at — send alert",
			issues: []domain.IssueExtended{
				{
					Issue: domain.Issue{
						ID:        3,
						ProjectID: 3,
						Level:     domain.IssueLevelError,
						CreatedAt: firstSeen,
						Status:    domain.IssueStatusUnresolved,
					},
					ResolvedAt: ptrTime(regressAt),
				},
			},
			expectedQueued: 1,
		},
		{
			name: "regress issue, last_notification_at is today — dont send alert",
			issues: []domain.IssueExtended{
				{
					Issue: domain.Issue{
						ID:                 4,
						ProjectID:          4,
						Level:              domain.IssueLevelError,
						CreatedAt:          firstSeen,
						Status:             domain.IssueStatusUnresolved,
						LastNotificationAt: ptrTime(today),
					},
					ResolvedAt: ptrTime(regressAt),
				},
			},
			expectedQueued: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			issuesRepo := &mockcontract.MockIssuesRepository{}
			notifQueueRepo := &mockcontract.MockNotificationsQueueRepository{}

			issuesRepo.EXPECT().ListUnresolved(mock.Anything).Return(tc.issues, tc.listErr)

			job := &AlertsDailyJob{
				issuesRepo:     issuesRepo,
				notifQueueRepo: notifQueueRepo,
				strategy:       newThroughOneStrategy(6),
			}

			if tc.issues != nil && tc.listErr == nil {
				for _, issue := range tc.issues {
					isNew := issue.ResolvedAt == nil
					shouldTrySend := false
					if issue.LastNotificationAt == nil {
						shouldTrySend = true
					} else {
						today := time.Now().Truncate(24 * time.Hour)
						lastSentAtDay := issue.LastNotificationAt.Truncate(24 * time.Hour)
						var firstDateDay time.Time
						if isNew {
							firstDateDay = issue.CreatedAt.Truncate(24 * time.Hour)
						} else {
							firstDateDay = *issue.ResolvedAt
						}
						shouldTrySend = job.strategy.Present(int(today.Sub(firstDateDay).Hours()/24)) &&
							firstDateDay != lastSentAtDay && lastSentAtDay != today
					}
					if shouldTrySend {
						notifQueueRepo.EXPECT().AddNotification(
							mock.Anything,
							issue.ProjectID,
							issue.ID,
							issue.Level,
							isNew,
							!isNew,
						).Return(tc.addNotifErr)
					}
				}
			}

			err := job.Run(context.Background())
			if tc.listErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			issuesRepo.AssertExpectations(t)
			notifQueueRepo.AssertExpectations(t)
		})
	}
}

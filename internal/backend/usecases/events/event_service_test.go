package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockIssueRepo := mockcontract.NewMockIssuesRepository(t)
	mockEventRepo := mockcontract.NewMockEventRepository(t)

	// Create service
	service := New(
		mockIssueRepo,
		mockEventRepo,
	)
	// Verify service was created correctly
	require.NotNil(t, service)
	require.Equal(t, mockIssueRepo, service.issueRepo)
	require.Equal(t, mockEventRepo, service.eventRepo)
}

func TestTimeseries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		filter        *domain.EventTimeseriesFilter
		setupMocks    func(mockEventRepo *mockcontract.MockEventRepository)
		expected      []domain.Timeseries
		expectedError bool
		errorContains string
	}{
		{
			name: "Basic timeseries",
			filter: &domain.EventTimeseriesFilter{
				Period: domain.Period{
					Interval:    24 * time.Hour,
					Granularity: time.Hour,
				},
				ProjectID: func() *domain.ProjectID {
					id := domain.ProjectID(1)

					return &id
				}(),
				Levels:  []domain.IssueLevel{domain.IssueLevelError},
				GroupBy: domain.EventTimeseriesGroupNone,
			},
			setupMocks: func(mockEventRepo *mockcontract.MockEventRepository) {
				mockEventRepo.EXPECT().Timeseries(mock.Anything, mock.AnythingOfType("*domain.EventTimeseriesFilter")).
					Return([]domain.Timeseries{
						{
							Name: "test",
							Period: domain.Period{
								Interval:    24 * time.Hour,
								Granularity: time.Hour,
							},
							Occurrences: []uint{10},
						},
					}, nil)
			},
			expected: []domain.Timeseries{
				{
					Name: "test",
					Period: domain.Period{
						Interval:    24 * time.Hour,
						Granularity: time.Hour,
					},
					Occurrences: []uint{10},
				},
			},
			expectedError: false,
		},
		{
			name: "Error from repository",
			filter: &domain.EventTimeseriesFilter{
				Period: domain.Period{
					Interval:    24 * time.Hour,
					Granularity: time.Hour,
				},
				ProjectID: func() *domain.ProjectID {
					id := domain.ProjectID(1)

					return &id
				}(),
				Levels:  []domain.IssueLevel{domain.IssueLevelError},
				GroupBy: domain.EventTimeseriesGroupNone,
			},
			setupMocks: func(mockEventRepo *mockcontract.MockEventRepository) {
				mockEventRepo.EXPECT().Timeseries(mock.Anything, mock.AnythingOfType("*domain.EventTimeseriesFilter")).
					Return(nil, errors.New("database error"))
			},
			expected:      nil,
			expectedError: true,
			errorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockIssueRepo := mockcontract.NewMockIssuesRepository(t)
			mockEventRepo := mockcontract.NewMockEventRepository(t)

			// Setup mocks
			tt.setupMocks(mockEventRepo)

			// Create service
			service := New(
				mockIssueRepo,
				mockEventRepo,
			)

			// Call the method
			result, err := service.Timeseries(context.Background(), tt.filter)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.expected), len(result))

				// Compare the timeseries
				for i := range tt.expected {
					require.Equal(t, tt.expected[i].Name, result[i].Name)
					require.Equal(t, tt.expected[i].Occurrences, result[i].Occurrences)
				}
			}
		})
	}
}

func TestIssueTimeseries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		filter        *domain.IssueEventsTimeseriesFilter
		setupMocks    func(mockIssueRepo *mockcontract.MockIssuesRepository, mockEventRepo *mockcontract.MockEventRepository)
		expected      []domain.Timeseries
		expectedError bool
		errorContains string
	}{
		{
			name: "Basic issue timeseries",
			filter: &domain.IssueEventsTimeseriesFilter{
				Period: domain.Period{
					Interval:    24 * time.Hour,
					Granularity: time.Hour,
				},
				ProjectID: 1,
				IssueID:   123,
				Levels:    []domain.IssueLevel{domain.IssueLevelError},
				GroupBy:   domain.EventTimeseriesGroupNone,
			},
			setupMocks: func(mockIssueRepo *mockcontract.MockIssuesRepository, mockEventRepo *mockcontract.MockEventRepository) {
				mockIssueRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(123)).
					Return(domain.Issue{
						ID:          123,
						Fingerprint: "test-fingerprint",
					}, nil)

				mockEventRepo.EXPECT().IssueTimeseries(mock.Anything, "test-fingerprint", mock.AnythingOfType("*domain.IssueEventsTimeseriesFilter")).
					Return([]domain.Timeseries{
						{
							Name: "test",
							Period: domain.Period{
								Interval:    24 * time.Hour,
								Granularity: time.Hour,
							},
							Occurrences: []uint{5},
						},
					}, nil)
			},
			expected: []domain.Timeseries{
				{
					Name: "test",
					Period: domain.Period{
						Interval:    24 * time.Hour,
						Granularity: time.Hour,
					},
					Occurrences: []uint{5},
				},
			},
			expectedError: false,
		},
		{
			name: "Error getting issue",
			filter: &domain.IssueEventsTimeseriesFilter{
				Period: domain.Period{
					Interval:    24 * time.Hour,
					Granularity: time.Hour,
				},
				ProjectID: 1,
				IssueID:   123,
				Levels:    []domain.IssueLevel{domain.IssueLevelError},
				GroupBy:   domain.EventTimeseriesGroupNone,
			},
			setupMocks: func(mockIssueRepo *mockcontract.MockIssuesRepository, mockEventRepo *mockcontract.MockEventRepository) {
				mockIssueRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(123)).
					Return(domain.Issue{}, errors.New("issue not found"))
			},
			expected:      nil,
			expectedError: true,
			errorContains: "get issue: issue not found",
		},
		{
			name: "Error getting timeseries",
			filter: &domain.IssueEventsTimeseriesFilter{
				Period: domain.Period{
					Interval:    24 * time.Hour,
					Granularity: time.Hour,
				},
				ProjectID: 1,
				IssueID:   123,
				Levels:    []domain.IssueLevel{domain.IssueLevelError},
				GroupBy:   domain.EventTimeseriesGroupNone,
			},
			setupMocks: func(mockIssueRepo *mockcontract.MockIssuesRepository, mockEventRepo *mockcontract.MockEventRepository) {
				mockIssueRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(123)).
					Return(domain.Issue{
						ID:          123,
						Fingerprint: "test-fingerprint",
					}, nil)

				mockEventRepo.EXPECT().IssueTimeseries(mock.Anything, "test-fingerprint", mock.AnythingOfType("*domain.IssueEventsTimeseriesFilter")).
					Return(nil, errors.New("database error"))
			},
			expected:      nil,
			expectedError: true,
			errorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockIssueRepo := mockcontract.NewMockIssuesRepository(t)
			mockEventRepo := mockcontract.NewMockEventRepository(t)

			// Setup mocks
			tt.setupMocks(mockIssueRepo, mockEventRepo)

			// Create service
			service := New(
				mockIssueRepo,
				mockEventRepo,
			)

			// Call the method
			result, err := service.IssueTimeseries(context.Background(), tt.filter)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.expected), len(result))

				// Compare the timeseries
				for i := range tt.expected {
					require.Equal(t, tt.expected[i].Name, result[i].Name)
					require.Equal(t, tt.expected[i].Occurrences, result[i].Occurrences)
				}
			}
		})
	}
}

package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func TestDomainUsersToAPI(t *testing.T) {
	testCases := []struct {
		name     string
		input    []domain.User
		expected []generatedapi.User
	}{
		{
			name:     "empty input",
			input:    []domain.User{},
			expected: []generatedapi.User{
				// empty output expected
			},
		},
		{
			name: "single user with nil last login",
			input: []domain.User{
				{
					ID:          1,
					Username:    "testuser",
					Email:       "testuser@example.com",
					IsSuperuser: true,
					IsActive:    true,
					CreatedAt:   time.Date(2025, 7, 1, 12, 0, 0, 0, time.UTC),
					LastLogin:   nil,
				},
			},
			expected: []generatedapi.User{
				{
					ID:          1,
					Username:    "testuser",
					Email:       "testuser@example.com",
					IsSuperuser: true,
					IsActive:    true,
					CreatedAt:   time.Date(2025, 7, 1, 12, 0, 0, 0, time.UTC),
					LastLogin:   generatedapi.OptDateTime{},
				},
			},
		},
		{
			name: "single user with non-nil last login",
			input: []domain.User{
				{
					ID:          2,
					Username:    "activeuser",
					Email:       "activeuser@example.com",
					IsSuperuser: false,
					IsActive:    true,
					CreatedAt:   time.Date(2024, 12, 15, 9, 30, 0, 0, time.UTC),
					LastLogin:   pointerToTime(time.Date(2025, 6, 25, 14, 20, 0, 0, time.UTC)),
				},
			},
			expected: []generatedapi.User{
				{
					ID:          2,
					Username:    "activeuser",
					Email:       "activeuser@example.com",
					IsSuperuser: false,
					IsActive:    true,
					CreatedAt:   time.Date(2024, 12, 15, 9, 30, 0, 0, time.UTC),
					LastLogin: generatedapi.OptDateTime{
						Value: time.Date(2025, 6, 25, 14, 20, 0, 0, time.UTC),
						Set:   true,
					},
				},
			},
		},
		{
			name: "multiple users with mixed last login states",
			input: []domain.User{
				{
					ID:          1,
					Username:    "user1",
					Email:       "user1@example.com",
					IsSuperuser: true,
					IsActive:    true,
					CreatedAt:   time.Date(2023, 3, 10, 8, 0, 0, 0, time.UTC),
					LastLogin:   nil,
				},
				{
					ID:          2,
					Username:    "user2",
					Email:       "user2@example.com",
					IsSuperuser: false,
					IsActive:    false,
					CreatedAt:   time.Date(2022, 1, 1, 15, 0, 0, 0, time.UTC),
					LastLogin:   pointerToTime(time.Date(2023, 5, 5, 10, 0, 0, 0, time.UTC)),
				},
			},
			expected: []generatedapi.User{
				{
					ID:          1,
					Username:    "user1",
					Email:       "user1@example.com",
					IsSuperuser: true,
					IsActive:    true,
					CreatedAt:   time.Date(2023, 3, 10, 8, 0, 0, 0, time.UTC),
					LastLogin:   generatedapi.OptDateTime{},
				},
				{
					ID:          2,
					Username:    "user2",
					Email:       "user2@example.com",
					IsSuperuser: false,
					IsActive:    false,
					CreatedAt:   time.Date(2022, 1, 1, 15, 0, 0, 0, time.UTC),
					LastLogin: generatedapi.OptDateTime{
						Value: time.Date(2023, 5, 5, 10, 0, 0, 0, time.UTC),
						Set:   true,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DomainUsersToAPI(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func pointerToTime(t time.Time) *time.Time {
	return &t
}

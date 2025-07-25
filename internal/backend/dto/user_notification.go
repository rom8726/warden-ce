package dto

import (
	"encoding/json"
	"fmt"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

//nolint:nestif // need refactor
func UserNotificationToDTO(notification domain.UserNotification) (generatedapi.UserNotification, error) {
	// Convert a domain type to an API type
	var apiType generatedapi.UserNotificationType
	switch notification.Type {
	case domain.UserNotificationTypeTeamAdded:
		apiType = generatedapi.UserNotificationTypeTeamAdded
	case domain.UserNotificationTypeTeamRemoved:
		apiType = generatedapi.UserNotificationTypeTeamRemoved
	case domain.UserNotificationTypeRoleChanged:
		apiType = generatedapi.UserNotificationTypeRoleChanged
	case domain.UserNotificationTypeIssueRegression:
		apiType = generatedapi.UserNotificationTypeIssueRegression
	default:
		apiType = generatedapi.UserNotificationTypeTeamAdded // fallback
	}

	// Convert content to map[string]jx.Raw
	contentMap := make(generatedapi.UserNotificationContent)
	if notification.Content != nil {
		var notifContent domain.UserNotificationContent
		if err := json.Unmarshal(notification.Content, &notifContent); err != nil {
			return generatedapi.UserNotification{}, fmt.Errorf("unmarshal notification content: %w", err)
		}

		var concrete any
		switch notification.Type {
		case domain.UserNotificationTypeTeamAdded:
			concrete = notifContent.TeamAdded
		case domain.UserNotificationTypeTeamRemoved:
			concrete = notifContent.TeamRemoved
		case domain.UserNotificationTypeRoleChanged:
			concrete = notifContent.RoleChanged
		case domain.UserNotificationTypeIssueRegression:
			concrete = notifContent.IssueRegression
		default:
			err := fmt.Errorf("unknown notification type: %s", notification.Type)

			return generatedapi.UserNotification{}, err
		}

		concreteContent, err := json.Marshal(concrete)
		if err != nil {
			return generatedapi.UserNotification{}, fmt.Errorf("marshal notification content: %w", err)
		}

		// Parse the JSON content and convert to map
		var contentData map[string]interface{}
		if err := json.Unmarshal(concreteContent, &contentData); err == nil {
			for key, value := range contentData {
				if jsonBytes, err := json.Marshal(value); err == nil {
					contentMap[key] = jsonBytes
				}
			}
		}
	}

	return generatedapi.UserNotification{
		ID:        uint(notification.ID),
		UserID:    uint(notification.UserID),
		Type:      apiType,
		Content:   contentMap,
		IsRead:    notification.IsRead,
		EmailSent: notification.EmailSent,
		CreatedAt: notification.CreatedAt,
		UpdatedAt: notification.UpdatedAt,
	}, nil
}

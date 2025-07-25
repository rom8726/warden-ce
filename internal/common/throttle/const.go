package throttle

import (
	"fmt"

	"github.com/rom8726/warden/internal/domain"
)

const (
	eventsCountMapKeyFormat = "project:%d:issue:event_counts"
	eventsIndexSetKeyFormat = "project:%d:issue:event_indexes"
)

func EventsCountMapKey(projectID domain.ProjectID) string {
	return fmt.Sprintf(eventsCountMapKeyFormat, projectID)
}

func EventsIndexSetKey(projectID domain.ProjectID) string {
	return fmt.Sprintf(eventsIndexSetKeyFormat, projectID)
}

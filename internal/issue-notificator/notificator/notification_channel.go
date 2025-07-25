package notificator

import (
	"context"
	"encoding/json"

	"github.com/rom8726/warden/internal/domain"
)

type Channel interface {
	Type() domain.NotificationType
	Send(
		ctx context.Context,
		issue *domain.Issue,
		project *domain.Project,
		config json.RawMessage,
		isRegress bool,
	) error
}

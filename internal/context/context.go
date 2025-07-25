//nolint:revive // it's ok here
package context

import (
	"context"
	"net/http"

	"github.com/rom8726/warden/internal/domain"
)

type contextKey string

const (
	ctxKeyProjectID contextKey = "project_id"
	ctxKeyUserID    contextKey = "user_id"
	ctxKeyIsSuper   contextKey = "is_superuser"
	ctxRawRequest   contextKey = "raw_request"
)

func WithProjectID(ctx context.Context, id domain.ProjectID) context.Context {
	return context.WithValue(ctx, ctxKeyProjectID, id)
}

func ProjectID(ctx context.Context) domain.ProjectID {
	return ctx.Value(ctxKeyProjectID).(domain.ProjectID)
}

func WithUserID(ctx context.Context, userID domain.UserID) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, userID)
}

func WithIsSuper(ctx context.Context, isSuper bool) context.Context {
	return context.WithValue(ctx, ctxKeyIsSuper, isSuper)
}

func IsSuper(ctx context.Context) bool {
	v, ok := ctx.Value(ctxKeyIsSuper).(bool)

	return ok && v
}

func UserID(ctx context.Context) domain.UserID {
	id, ok := ctx.Value(ctxKeyUserID).(domain.UserID)
	if !ok {
		return 0
	}

	return id
}

func WithRawRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, ctxRawRequest, req)
}

func RawRequest(ctx context.Context) *http.Request {
	return ctx.Value(ctxRawRequest).(*http.Request)
}

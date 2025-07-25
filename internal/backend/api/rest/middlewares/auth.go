package middlewares

import (
	"net/http"
	"strings"

	"github.com/rom8726/warden/internal/backend/contract"
	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
)

// AuthMiddleware extracts the user ID from the request and sets it in the context.
func AuthMiddleware(tokenizer contract.Tokenizer, usersSrv contract.UsersUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// Extract the Authorization header
			authHeader := request.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				// No auth header or not a bearer token, just pass through
				next.ServeHTTP(writer, request)

				return
			}

			// Extract the token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Verify the token and get the user ID
			claims, err := tokenizer.VerifyToken(token, domain.TokenTypeAccess)
			if err != nil {
				// Invalid token, pass through
				next.ServeHTTP(writer, request)

				return
			}

			// Get the user
			user, err := usersSrv.GetByID(request.Context(), domain.UserID(claims.UserID))
			if err != nil {
				// User isn't found, pass through
				next.ServeHTTP(writer, request)

				return
			}

			// Set the user ID and superuser flag in the context
			ctx := wardencontext.WithUserID(request.Context(), user.ID)
			ctx = wardencontext.WithIsSuper(ctx, user.IsSuperuser)

			// Continue with the modified context
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

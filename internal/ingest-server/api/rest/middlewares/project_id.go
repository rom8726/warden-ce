package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
)

func WithProjectID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		vars := request.URL.Path
		parts := strings.Split(vars, "/")

		if len(parts) >= 5 {
			if !(parts[1] == "api" && (parts[3] == "envelope" || parts[3] == "store")) {
				next.ServeHTTP(writer, request)

				return
			}

			projectIDStr := parts[2]
			projectIDUint, err := strconv.ParseUint(projectIDStr, 10, 64)
			if err != nil {
				http.Error(writer, "Invalid project ID", http.StatusBadRequest)

				return
			}

			ctx := context.WithProjectID(request.Context(), domain.ProjectID(projectIDUint))
			next.ServeHTTP(writer, request.WithContext(ctx))

			return
		}

		next.ServeHTTP(writer, request)
	})
}

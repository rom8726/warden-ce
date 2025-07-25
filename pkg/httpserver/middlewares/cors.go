package middlewares

import (
	"net/http"
)

func CORSMdw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if req.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(writer, req)
	})
}

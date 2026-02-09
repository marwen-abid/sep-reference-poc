package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/stellar/sep-reference/reference/go/sep10"
)

type contextKey string

const accountKey contextKey = "sep10-account"

func AccountFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(accountKey).(string); ok {
		return v
	}
	return ""
}

func SEP10Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := strings.TrimSpace(r.Header.Get("Authorization"))
			if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				writeJSONError(w, http.StatusUnauthorized, "missing bearer token")
				return
			}
			token := strings.TrimSpace(auth[7:])
			claims, err := sep10.VerifyToken(token, jwtSecret)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "invalid bearer token")
				return
			}
			ctx := context.WithValue(r.Context(), accountKey, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

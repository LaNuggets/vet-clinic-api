package authentication

import (
	"context"
	"net/http"
)

// Middleware to secure routes with a JWT
func AuthMiddleware(secret string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Check Header content
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			// Check token validity
			claims, err := ParseTokenClaims(secret, authHeader)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "email", claims["email"])
			ctx = context.WithValue(ctx, "role", claims["role"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Check if the required Role is match by the user token
func RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value("role").(string)
			if !ok || role != requiredRole {
				http.Error(w, "Forbidden: insufficient privileges", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

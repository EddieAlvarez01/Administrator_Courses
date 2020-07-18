package middlewares

import (
	"context"
	"github.com/EddieAlvarez01/administrator_courses/authorization"
	"github.com/EddieAlvarez01/administrator_courses/models"
	"net/http"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-access-token")
		if tokenString == "" {
			models.NewResponseJSON(w, http.StatusUnauthorized, "No token provided", nil)
			return
		}
		payload, err := authorization.VerifyToken(tokenString)
		if err != nil {
			models.NewResponseJSON(w, http.StatusUnauthorized, "Invalid token", nil)
			return
		}
		ctx := context.WithValue(r.Context(), "payload", payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleAdministrator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := r.Context().Value("payload").(models.Payload)
		for _, role := range payload.Role {
			if role == "ADMIN" {
				next.ServeHTTP(w, r)
				return
			}
		}
		models.NewResponseJSON(w, http.StatusForbidden, "Access denied", nil)
	})
}

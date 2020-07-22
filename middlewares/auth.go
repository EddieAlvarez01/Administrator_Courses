package middlewares

import (
	"context"
	"fmt"
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

//VERIFY USER ROLE | FLAG 0 ADMIN, 1 STUDENT, 2 PROFESSOR
func PersonRole(next http.Handler, flag uint8) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := r.Context().Value("payload").(models.Payload)
		for _, role := range payload.Role {
			switch flag {
			case 0:
				if role == "ADMIN" {
					next.ServeHTTP(w, r)
					return
				}
			case 1:
				if role == "STUDENT" {
					next.ServeHTTP(w, r)
					return
				}
			case 2:
				if role == "PROFESSOR"{
					next.ServeHTTP(w, r)
					return
				}
			default:
				fmt.Println("Invalid flag")
				models.NewResponseJSON(w, http.StatusInternalServerError, "Server error", nil)
				return
			}
		}
		models.NewResponseJSON(w, http.StatusForbidden, "Access denied", nil)
	})
}

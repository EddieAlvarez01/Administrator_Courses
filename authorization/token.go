package authorization

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EddieAlvarez01/administrator_courses/models"
	"github.com/dgrijalva/jwt-go"
)

//GenerateToken RETURN A JWT TOKEN
func GenerateToken(payload models.Person) (string, error) {
	payloadJWT := models.Payload{
		ID:   payload.ID.Hex(),
		Role: payload.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 2).Unix(),
			Issuer:    "Eddie Alvarez",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payloadJWT)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return tokenString, nil
}

//VerifyToken MIDDLEWARE VERIFY TOKEN
func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-access-token")
		if tokenString == "" {
			models.NewResponseJSON(w, http.StatusUnauthorized, "Token no provided", nil)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRS256); !ok {
				return nil, 
			}
		})
	})
}

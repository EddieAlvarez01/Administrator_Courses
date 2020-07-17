package authorization

import (
	"errors"
	"fmt"
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
func VerifyToken(tokenString string) (models.Payload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})
	if err != nil {
		return models.Payload{}, err
	}
	if payload, ok := token.Claims.(models.Payload); ok && token.Valid {
		return payload, nil
	}
	return models.Payload{}, errors.New("Invalid token")
}

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
	token, err := jwt.ParseWithClaims(tokenString, &models.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return models.Payload{}, err
	}
	payload, ok := token.Claims.(*models.Payload)
	if !token.Valid {
		fmt.Println("Invalid token")
		return models.Payload{}, errors.New("Invalid token")
	}
	if !ok {
		fmt.Println("Error getting the payload")
		return models.Payload{}, errors.New("Error getting the payload")
	}
	return *payload, nil
}

package models

import "github.com/dgrijalva/jwt-go"

//Payload .
type Payload struct {
	ID   string   `json:"id"`
	Role []string `json:"role"`
	jwt.StandardClaims
}

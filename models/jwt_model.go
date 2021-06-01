package model

import "github.com/dgrijalva/jwt-go"

//JwtKey -> secret key
var JwtKey = []byte("super_secret_phrase_key")

// Credentials ->
type Credentials struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

// Claims ->
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

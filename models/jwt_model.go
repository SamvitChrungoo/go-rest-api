package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-rest-api/utils"
)

//JwtKey -> secret key
var JwtKey = []byte(utils.GetEnvironmentVariable("JWT_SECRET_PHASE"))

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

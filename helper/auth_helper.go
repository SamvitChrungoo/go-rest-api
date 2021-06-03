package helper

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	model "github.com/go-rest-api/models"
)

//CheckAuth -> returns ture or false based on Auth status
func CheckAuth(request *http.Request) (bool, string, *model.Claims) {
	cookie, err := request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, "Unauthorized Access :( ", &model.Claims{}
		}
		return false, "Bad Request !!", &model.Claims{}
	}
	tokenStr := cookie.Value
	claims := &model.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, "Unauthorized Access :( ", &model.Claims{}
		}
		return false, "Unauthorized Access :( ", &model.Claims{}
	}
	if !tkn.Valid {
		return false, "Unauthorized Access :( ", &model.Claims{}
	}
	return true, "Success", claims
}

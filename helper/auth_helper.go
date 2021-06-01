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
			return false, "StatusUnauthorized", &model.Claims{}
		}
		return false, "StatusBadRequest", &model.Claims{}
	}

	tokenStr := cookie.Value
	claims := &model.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, "StatusUnauthorized", &model.Claims{}
		}
		return false, "StatusUnauthorized", &model.Claims{}
	}
	if !tkn.Valid {
		return false, "StatusUnauthorized", &model.Claims{}
	}
	return true, "Success", claims
}

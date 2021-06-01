package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	model "github.com/go-rest-api/models"
	"github.com/go-rest-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

// LoginHandler ->
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	var credentials model.Credentials
	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	var expextedCredentials model.Credentials
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client = utils.ConnectDB()
	userCollection := client.Database("test_db").Collection("users")
	result := userCollection.FindOne(ctx, bson.M{"username": credentials.Username}).Decode(&expextedCredentials)
	if result != nil {
		response.WriteHeader(http.StatusInternalServerError)
	}
	if credentials.Password != expextedCredentials.Password {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"message": "Unauthorized Access !!"}`))
		return
	}
	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &model.Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(model.JwtKey)
	fmt.Println(tokenString)
	http.SetCookie(response,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	response.Write([]byte(`{"message": "Successfully logged in"}`))
}

// LogOutHandler ->
func LogOutHandler(response http.ResponseWriter, request *http.Request) {
	http.SetCookie(response,
		&http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Unix(0, 0),
		})
	response.Write([]byte(`{"message": "Successfully logged out !!"}`))
}

// SignUpHandler ->
func SignUpHandler(response http.ResponseWriter, request *http.Request) {
	var credentials model.Credentials
	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client = utils.ConnectDB()
	userCollection := client.Database("test_db").Collection("users")
	result, _ := userCollection.InsertOne(ctx, credentials)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(response).Encode(result)
}

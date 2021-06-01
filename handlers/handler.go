package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-rest-api/helper"
	model "github.com/go-rest-api/models"
	"github.com/go-rest-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	defer client.Disconnect(ctx)
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
	defer client.Disconnect(ctx)
	userCollection := client.Database("test_db").Collection("users")
	result, _ := userCollection.InsertOne(ctx, credentials)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(response).Encode(result)
}

// AddMovieEndpoint -> to add movie to DB
func AddMovieEndpoint(response http.ResponseWriter, request *http.Request) {
	isAuthenticated, authStatus, _ := helper.CheckAuth(request)
	if isAuthenticated == false {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"message": "` + authStatus + `"}`))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client = utils.ConnectDB()
	defer client.Disconnect(ctx)
	response.Header().Add("content-type", "application/json")
	var movie model.Movie
	json.NewDecoder(request.Body).Decode(&movie)
	movieCollection := client.Database("test_db").Collection("movies")
	defer cancel()
	result, _ := movieCollection.InsertOne(ctx, movie)
	json.NewEncoder(response).Encode(result)
}

// DeleteMovieEndpoint -> to add movie to DB
func DeleteMovieEndpoint(response http.ResponseWriter, request *http.Request) {
	isAuthenticated, authStatus, _ := helper.CheckAuth(request)
	if isAuthenticated == false {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"message": "` + authStatus + `"}`))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client = utils.ConnectDB()
	defer client.Disconnect(ctx)
	response.Header().Add("content-type", "application/json")
	movieCollection := client.Database("test_db").Collection("movies")
	defer cancel()
	deletionID := request.FormValue("id")
	if deletionID == "" {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Bad request (check id)"}`))
		return
	}
	id, _ := primitive.ObjectIDFromHex(deletionID)
	result, _ := movieCollection.DeleteOne(ctx, bson.M{"_id": id})
	json.NewEncoder(response).Encode(result)
}

// UpdateMovieEndpoint -> to add movie to DB
func UpdateMovieEndpoint(response http.ResponseWriter, request *http.Request) {
	isAuthenticated, authStatus, _ := helper.CheckAuth(request)
	if isAuthenticated == false {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"message": "` + authStatus + `"}`))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client = utils.ConnectDB()
	defer client.Disconnect(ctx)
	response.Header().Add("content-type", "application/json")
	var movie model.Movie
	json.NewDecoder(request.Body).Decode(&movie)
	movieCollection := client.Database("test_db").Collection("movies")
	defer cancel()
	updationID := request.FormValue("id")
	id, _ := primitive.ObjectIDFromHex(updationID)
	update := bson.M{"$set": bson.M{
		"title":         movie.Title,
		"imdbRating":    movie.ImdbRating,
		"originalTitle": movie.OriginalTitle,
		"posterurl":     movie.PosterURL,
		"storyline":     movie.Storyline,
		"year":          movie.Year}}
	result, _ := movieCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	json.NewEncoder(response).Encode(result)
}

// GetMovieEndpoint -> to get one or all movies from DB
func GetMovieEndpoint(response http.ResponseWriter, request *http.Request) {
	isAuthenticated, authStatus, _ := helper.CheckAuth(request)
	if isAuthenticated == false {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"message": "` + authStatus + `"}`))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client = utils.ConnectDB()
	defer client.Disconnect(ctx)
	response.Header().Add("content-type", "application/json")
	movieCollection := client.Database("test_db").Collection("movies")
	allIDs := request.URL.Query()["id"]
	var movies []model.Movie
	if len(allIDs) > 0 {
		var slice = make([]primitive.ObjectID, len(allIDs))
		for i := 0; i < len(allIDs); i++ {
			id, _ := primitive.ObjectIDFromHex(allIDs[i])
			slice = append(slice, id)
		}
		output := slice[len(allIDs):]
		query := bson.M{"_id": bson.M{"$in": output}}
		result, err := movieCollection.Find(ctx, query)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		}
		defer result.Close(ctx)
		for result.Next(ctx) {
			var movie model.Movie
			result.Decode(&movie)
			movies = append(movies, movie)
		}
		if err := result.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}
		json.NewEncoder(response).Encode(movies)

	} else {
		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var movie model.Movie
			cursor.Decode(&movie)
			movies = append(movies, movie)
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}
		json.NewEncoder(response).Encode(movies)
	}

}

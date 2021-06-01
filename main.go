package main

import (
	"context"
	"encoding/json"
	"net/http"

	"time"

	handler "github.com/go-rest-api/handlers"
	helper "github.com/go-rest-api/helper"
	model "github.com/go-rest-api/models"
	utils "github.com/go-rest-api/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func main() {
	client = utils.ConnectDB()
	defer utils.DisconnectDB()
	router := mux.NewRouter()
	router.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	router.HandleFunc("/signup", handler.SignUpHandler).Methods("POST")
	router.HandleFunc("/movie", AddMovieEndpoint).Methods("POST")
	router.HandleFunc("/movies", GetMovieEndpoint).Methods("GET")
	http.ListenAndServe(":8080", router)
}

// AddMovieEndpoint -> to add movie to DB
func AddMovieEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var movie model.Movie
	json.NewDecoder(request.Body).Decode(&movie)
	movieCollection := client.Database("test_db").Collection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := movieCollection.InsertOne(ctx, movie)
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
	response.Header().Add("content-type", "application/json")
	movieCollection := client.Database("test_db").Collection("movies")
	requestedID := request.FormValue("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if requestedID != "" {
		var movie model.Movie
		id, _ := primitive.ObjectIDFromHex(requestedID)
		result := movieCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&movie)
		if result != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": "` + result.Error() + `"}`))
		} else {
			json.NewEncoder(response).Encode(movie)
		}
	} else {
		var movies []model.Movie
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

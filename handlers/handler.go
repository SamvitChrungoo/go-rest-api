package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-rest-api/helper"
	model "github.com/go-rest-api/models"
	"github.com/go-rest-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoginHandler ->
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	var credentials model.Credentials
	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		response.Header().Add("content-type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusBadRequest, Message: "Bad Request :-("})
		response.Write(errorResponse)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var expextedCredentials model.Credentials
	userCollection := utils.Client.Database(utils.GetEnvironmentVariable("DB_NAME")).Collection("users")
	result := userCollection.FindOne(ctx, bson.M{"username": credentials.Username})
	decodeError := result.Decode(&expextedCredentials)
	if decodeError != nil {
		response.Header().Add("content-type", "application/json")
		errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusInternalServerError, Message: "No user found :-("})
		response.WriteHeader(http.StatusUnauthorized)
		response.Write(errorResponse)
		return
	}

	if credentials.Password != expextedCredentials.Password {
		response.Header().Add("content-type", "application/json")
		errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusUnauthorized, Message: "Unauthorized Access - incorrect password :-("})
		response.WriteHeader(http.StatusUnauthorized)
		response.Write(errorResponse)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 10)
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

	response.Header().Add("content-type", "application/json")
	successResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusOK, Message: "Logged in successfully :-) "})
	response.WriteHeader(http.StatusOK)
	response.Write(successResponse)
}

// LogOutHandler ->
func LogOutHandler(response http.ResponseWriter, request *http.Request) {
	http.SetCookie(response,
		&http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Unix(0, 0),
		})

	response.Header().Add("content-type", "application/json")
	successResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusOK, Message: "Logged in successfully :-) "})
	response.WriteHeader(http.StatusOK)
	response.Write(successResponse)
}

// SignUpHandler ->
func SignUpHandler(response http.ResponseWriter, request *http.Request) {
	var credentials model.Credentials
	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		response.Header().Add("content-type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusBadRequest, Message: "Bad Request :-("})
		response.Write(errorResponse)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := utils.Client.Database(utils.GetEnvironmentVariable("DB_NAME")).Collection("users")
	result, _ := userCollection.InsertOne(ctx, credentials)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	}
	response.Header().Add("content-type", "application/json")
	json.NewEncoder(response).Encode(result)
}

// AddMovieEndpoint -> to add movie to DB
func AddMovieEndpoint(response http.ResponseWriter, request *http.Request) {
	isAuthenticated, authStatus, _ := helper.CheckAuth(request)
	if isAuthenticated == false {
		response.Header().Add("content-type", "application/json")
		errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusUnauthorized, Message: authStatus})
		response.WriteHeader(http.StatusUnauthorized)
		response.Write(errorResponse)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response.Header().Add("content-type", "application/json")
	var movie model.Movie
	json.NewDecoder(request.Body).Decode(&movie)
	movieCollection := utils.Client.Database(utils.GetEnvironmentVariable("DB_NAME")).Collection("movies")
	result, _ := movieCollection.InsertOne(ctx, movie)
	json.NewEncoder(response).Encode(result)
}

// DeleteMovieEndpoint -> to add movie to DB
func DeleteMovieEndpoint(response http.ResponseWriter, request *http.Request) {
	isAuthenticated, authStatus, _ := helper.CheckAuth(request)
	if isAuthenticated == false {
		response.Header().Add("content-type", "application/json")
		errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusUnauthorized, Message: authStatus})
		response.WriteHeader(http.StatusUnauthorized)
		response.Write(errorResponse)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response.Header().Add("content-type", "application/json")
	movieCollection := utils.Client.Database(utils.GetEnvironmentVariable("DB_NAME")).Collection("movies")
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
		response.Header().Add("content-type", "application/json")
		errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusUnauthorized, Message: authStatus})
		response.WriteHeader(http.StatusUnauthorized)
		response.Write(errorResponse)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response.Header().Add("content-type", "application/json")
	var movie model.Movie
	json.NewDecoder(request.Body).Decode(&movie)
	movieCollection := utils.Client.Database(utils.GetEnvironmentVariable("DB_NAME")).Collection("movies")
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
		response.Header().Add("content-type", "application/json")
		errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusUnauthorized, Message: authStatus})
		response.WriteHeader(http.StatusUnauthorized)
		response.Write(errorResponse)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response.Header().Add("content-type", "application/json")
	movieCollection := utils.Client.Database(utils.GetEnvironmentVariable("DB_NAME")).Collection("movies")
	allIDs := request.URL.Query()["id"]
	sortBy := request.FormValue("sortBy")
	sortOrder := request.FormValue("sortOrder")
	filter := request.FormValue("filter")
	min := request.FormValue("min")
	max := request.FormValue("max")
	var movies []model.Movie

	if len(allIDs) > 0 {
		var objectIDs = make([]primitive.ObjectID, len(allIDs))
		for i := 0; i < len(allIDs); i++ {
			id, _ := primitive.ObjectIDFromHex(allIDs[i])
			objectIDs[i] = id
		}
		query := bson.M{"_id": bson.M{"$in": objectIDs}}
		result, err := movieCollection.Find(ctx, query)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Header().Add("content-type", "application/json")
			errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusInternalServerError, Message: "Internal Server Error :-)"})
			response.Write(errorResponse)
			return
		}
		defer result.Close(ctx)
		for result.Next(ctx) {
			var movie model.Movie
			result.Decode(&movie)
			movies = append(movies, movie)
		}
		if err := result.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Header().Add("content-type", "application/json")
			errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusInternalServerError, Message: "Internal Server Error :-)"})
			response.Write(errorResponse)
			return
		}
	} else {
		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Header().Add("content-type", "application/json")
			errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusInternalServerError, Message: "Internal Server Error :-)"})
			response.Write(errorResponse)
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var movie model.Movie
			cursor.Decode(&movie)
			movies = append(movies, movie)
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Header().Add("content-type", "application/json")
			errorResponse, _ := json.Marshal(model.ErrorResponse{Code: http.StatusInternalServerError, Message: "Internal Server Error :-)"})
			response.Write(errorResponse)
			return
		}
	}
	var filteredResults = make([]model.Movie, 0)

	if filter == "rating" || filter == "year" {
		var minRating float64
		var maxRating float64
		var minYear int
		var maxYear int
		if filter == "rating" {
			if min == "" {
				minRating = 0.0
			} else {
				minRating, _ = strconv.ParseFloat(min, 64)
			}
			if max == "" {
				maxRating = 10.0
			} else {
				maxRating, _ = strconv.ParseFloat(max, 64)
			}
		} else {
			if min == "" {
				minYear = 1800
			} else {
				minYear, _ = strconv.Atoi(min)
			}
			if max == "" {
				maxYear = time.Now().Year()
			} else {
				maxYear, _ = strconv.Atoi(max)
			}
		}

		for i := 0; i < len(movies); i++ {
			if filter == "rating" {
				if movies[i].ImdbRating >= minRating && movies[i].ImdbRating <= maxRating {
					filteredResults = append(filteredResults, movies[i])
				}
			} else {
				if movies[i].Year >= minYear && movies[i].Year <= maxYear {
					filteredResults = append(filteredResults, movies[i])
				}
			}
		}
	}
	if len(filteredResults) == 0 {
		for i := 0; i < len(movies); i++ {
			filteredResults = append(filteredResults, movies[i])
		}
	}
	if sortBy == "rating" || sortBy == "year" {
		if sortOrder == "ascending" {
			if sortBy == "rating" {
				sort.Slice(filteredResults, func(i, j int) bool {
					return filteredResults[i].ImdbRating < filteredResults[j].ImdbRating
				})
			} else {
				sort.Slice(filteredResults, func(i, j int) bool {
					return filteredResults[i].Year < filteredResults[j].Year
				})
			}
		} else if sortOrder == "descending" || sortOrder == "" {
			if sortBy == "rating" {
				sort.Slice(filteredResults, func(i, j int) bool {
					return filteredResults[i].ImdbRating > filteredResults[j].ImdbRating
				})
			} else {
				sort.Slice(filteredResults, func(i, j int) bool {
					return filteredResults[i].Year > filteredResults[j].Year
				})
			}
		}
	}
	json.NewEncoder(response).Encode(filteredResults)
}

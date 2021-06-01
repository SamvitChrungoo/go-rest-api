package main

import (
	"net/http"

	handler "github.com/go-rest-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", handler.LogOutHandler).Methods("GET")
	router.HandleFunc("/signup", handler.SignUpHandler).Methods("POST")
	router.HandleFunc("/movie", handler.AddMovieEndpoint).Methods("POST")
	router.HandleFunc("/movie", handler.UpdateMovieEndpoint).Methods("PUT")
	router.HandleFunc("/movie", handler.DeleteMovieEndpoint).Methods("DELETE")
	router.HandleFunc("/movies", handler.GetMovieEndpoint).Methods("GET")
	http.ListenAndServe(":8080", router)
}

package main

import (
	"net/http"

	handler "github.com/go-rest-api/handlers"
	utils "github.com/go-rest-api/utils"

	"github.com/gorilla/mux"
)

func main() {
	utils.SetEnvironmentVariables()
	utils.ConnectDB()
	defer close()
	router := mux.NewRouter()
	router.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", handler.LogOutHandler).Methods("GET")
	router.HandleFunc("/signup", handler.SignUpHandler).Methods("POST")
	router.HandleFunc("/movie", handler.AddMovieEndpoint).Methods("POST")
	router.HandleFunc("/movie", handler.UpdateMovieEndpoint).Methods("PUT")
	router.HandleFunc("/movie", handler.DeleteMovieEndpoint).Methods("DELETE")
	router.HandleFunc("/movies", handler.GetMovieEndpoint).Methods("GET")
	portNo := utils.GetEnvironmentVariable("PORT_NO")
	http.ListenAndServe(":"+portNo, router)
}

func close() {
	utils.UnsetEnvironmentVariables()
	utils.DisconnectDB()
}

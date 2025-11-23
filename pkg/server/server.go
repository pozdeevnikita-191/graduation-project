package server

import (
	"main/pkg/api"
	"net/http"
	"os"
)

//Run starts the server
func Run() error {

	port := os.Getenv("TODO_PORT")
	if port == ""{
		port = ":7540"
	}

	api.Init()
	return http.ListenAndServe(port, nil)
}


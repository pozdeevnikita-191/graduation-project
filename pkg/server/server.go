package server

import (
	"main/pkg/api"
	"net/http"
)

//Run starts the server
func Run() error {
	api.Init()
	return http.ListenAndServe(":7540", nil)
}


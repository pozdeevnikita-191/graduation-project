package main

import (
	"fmt"
	"log"
	"main/pkg/db"
	"main/pkg/server"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	fmt.Println("Сервер запущен")
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal("failed to initialize the DB package:", err)
	}

	http.Handle("/", http.FileServer(http.Dir("web")))
	if err := server.Run(); err != nil {
		log.Fatal(err)
	} 
}



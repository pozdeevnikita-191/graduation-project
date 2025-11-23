package main

import (
	"fmt"
	"log"
	"main/pkg/db"
	"main/pkg/server"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if port == ""{
		port = ":7540"
	}
	
	fmt.Println("Server started. Port" + port)

	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal("failed to initialize the DB package:", err)
	}

	defer db.DB.Close()  

	http.Handle("/", http.FileServer(http.Dir("web")))
	if err := server.Run(); err != nil {
		log.Fatal(err)
	} 
}



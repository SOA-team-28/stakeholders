package main

import (
	"database-example/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func startServer() {
	database := db.InitDB()
	if database == nil {
		log.Fatal("FAILED TO CONNECT TO DB")
	}

	router := mux.NewRouter().StrictSlash(true)

	//dodati handlere

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8086", router))
}

func main() {

	startServer()
}

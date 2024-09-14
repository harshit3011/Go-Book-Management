package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshit3011/book-management-system/database"
	"github.com/harshit3011/book-management-system/routes"
)

func main() {
	database.InitDB()

	router:=mux.NewRouter()

	routes.RegisterUserRoutes(router)
	routes.RegisterBookRoutes(router)

	log.Fatal(http.ListenAndServe(":8080",router))
}
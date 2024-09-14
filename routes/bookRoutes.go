package routes

import (
	"github.com/gorilla/mux"
	"github.com/harshit3011/book-management-system/controllers"
	"github.com/harshit3011/book-management-system/middleware"
)

func RegisterBookRoutes(router *mux.Router){
	router.HandleFunc("/books", middleware.AuthMiddleware(controllers.AddBook)).Methods("POST")

	router.HandleFunc("/books", middleware.AuthMiddleware(controllers.GetBooks)).Methods("GET")

	router.HandleFunc("/books/{id}", middleware.AuthMiddleware(controllers.DeleteBook)).Methods("DELETE")
}
package routes

import (
	"github.com/gorilla/mux"
	"github.com/harshit3011/book-management-system/controllers"
)

func RegisterUserRoutes(router *mux.Router){
	router.HandleFunc("/signup",controllers.SignUpUser).Methods("POST")
	router.HandleFunc("/login",controllers.LoginUser).Methods("POST")
}
package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/harshit3011/book-management-system/database"
	"github.com/harshit3011/book-management-system/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



func SignUpUser(w http.ResponseWriter, r *http.Request){
	var user database.User

	err:= json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w,"Invalid Input", http.StatusBadRequest)
		return
	}
	err= user.HashPassword(user.Password)
	if err != nil {
		http.Error(w,"Error while Hashing the password",http.StatusInternalServerError)
		return
	}
	var users *mongo.Collection = database.Client.Database("bookshelf").Collection("users")
	user.Books = []primitive.ObjectID{}
	result,err:=users.InsertOne(context.TODO(),user)
	if err != nil {
		http.Error(w,"Error while saving the user",http.StatusInternalServerError)
		return
	}
	token,err:=utils.GenerateToken(result.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		http.Error(w,"Error while generating the token",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{"token":token})
}

func LoginUser(w http.ResponseWriter, r *http.Request){
	var newUser database.User
	err:= json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w,"Invalid input",http.StatusBadRequest)
		return
	}
	email:=newUser.Email
	password := newUser.Password

    var user database.User
	var users *mongo.Collection = database.Client.Database("bookshelf").Collection("users")

	err= users.FindOne(context.TODO(),bson.M{"email":email}).Decode(&user)
	if err != nil {
		http.Error(w,"Invalid email or password",http.StatusUnauthorized)
		return
	}
	err=user.CheckPassword(password)
	if err != nil {
		http.Error(w,"Invalid email or password",http.StatusUnauthorized)
		return
	}
	token,err:= utils.GenerateToken(user.Id.Hex())
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})

}
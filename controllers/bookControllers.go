package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/harshit3011/book-management-system/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func AddBook(w http.ResponseWriter, r *http.Request){
	userId := r.Context().Value("user_id").(string)
    objectUserId, err := primitive.ObjectIDFromHex(userId)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var book database.Book
    err = json.NewDecoder(r.Body).Decode(&book)
    if err != nil {
        http.Error(w, "Invalid book credentials", http.StatusBadRequest)
        return
    }
    book.UserID = objectUserId
    book.Id = primitive.NewObjectID()

    ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
    defer cancel()

    books := database.Client.Database("bookshelf").Collection("books")
    userCollection := database.Client.Database("bookshelf").Collection("users")

    _, err = books.InsertOne(ctx, book)
    if err != nil {
        http.Error(w, "Error while saving the book", http.StatusInternalServerError)
        return
    }

    // Check if the user exists
    var user database.User
    err = userCollection.FindOne(ctx, bson.M{"_id": objectUserId}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Error finding user", http.StatusInternalServerError)
        return
    }
	fmt.Println(user, objectUserId, book.Id)

    result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objectUserId}, bson.M{"$push": bson.M{"books": book.Id}})
    if err != nil {
		fmt.Println(err)
        http.Error(w, "Failed to link book to user", http.StatusInternalServerError)
        return
    }

    if result.MatchedCount == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Book added successfully"})
	}

func GetBooks(w http.ResponseWriter, r *http.Request){
	userId := r.Context().Value("user_id").(string)
    objectUserId, err := primitive.ObjectIDFromHex(userId)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var myBooks []database.Book
    filter := bson.M{"user_id": objectUserId}
    cursor, err := database.Client.Database("bookshelf").Collection("books").Find(ctx, filter)
    if err != nil {
        http.Error(w, "Error fetching books", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    err = cursor.All(ctx, &myBooks)
    if err != nil {
        http.Error(w, "Error processing books", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(myBooks)

}
func DeleteBook(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value("user_id").(string)
	objectUserId, _ := primitive.ObjectIDFromHex(userID)

	params:=mux.Vars(r)
	bookId,err:=primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w,"Invalid Book ID",http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	var books *mongo.Collection = database.Client.Database("bookshelf").Collection("books")
var userCollection *mongo.Collection = database.Client.Database("bookshelf").Collection("users")
	filter:=bson.M{"_id":bookId,"user_id": objectUserId}
	result:=books.FindOne(ctx,filter)
	if result.Err()!=nil{
		
		http.Error(w, "No book found or you don't have permission to delete it", http.StatusNotFound)
		return
	}
	_,err=userCollection.UpdateOne(ctx,bson.M{"_id": objectUserId}, bson.M{"$pull": bson.M{"books": bookId}})
	if err!=nil {
		http.Error(w, "Book could not be deleted", http.StatusInternalServerError)
		return
	}

	_,err=books.DeleteOne(ctx,bson.M{"_id":bookId})
	if err != nil {
		http.Error(w, "Error deleting the book", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted successfully"})
}
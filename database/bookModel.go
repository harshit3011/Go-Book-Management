package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title string `bson:"title" json:"title"`
	Author string `bson:"author" json:"author"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
}

package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var Client *mongo.Client

func Connect() *mongo.Client{
	err:= godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error while fetching the environment variables%s",err)
	}
	MONGOURI:=os.Getenv("MONGO_URI")
	if MONGOURI == "" {
		log.Fatalf("MONGO_URI environment variable is not set")
	}
	clientOptions:=options.Client().ApplyURI(MONGOURI)
	client,err:=mongo.Connect(context.Background(),clientOptions)
	if err != nil {
		log.Fatalf("%s",err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

func InitDB(){
	Client=Connect()
	if Client == nil {
		log.Fatalf("Failed to initialize MongoDB client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
}
package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/bson"
	"lambda-push-go/core"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DB_NAME")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx = context.WithValue(ctx, core.DbURL, dbUrl)
	db, err := core.ConfigDB(ctx, dbName)
	if err != nil {
		log.Fatalf("database configuration failed: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	var notification core.Notification
	err = db.Collection("notifications").FindOne(ctx, bson.D{}).Decode(&notification)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("notification: %+v", notification)
}

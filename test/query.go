package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"log"
	"os"

	"github.com/mongodb/mongo-go-driver/bson"
	"lambda-push-go/core"
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
	var notificationAccount core.NotificationAccount

	id, _ := primitive.ObjectIDFromHex("5c9dea2e6ffa8676f90163d6")
	err = db.Collection("notifications").FindOne(ctx, bson.M{ "_id": id }).Decode(&notification)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("notification:", notification)


	err = db.Collection("notificationaccounts").FindOne(ctx, bson.M{"siteId": notification.SiteID}).Decode(&notificationAccount)

	if err != nil {
		log.Fatal(err)
	}
}

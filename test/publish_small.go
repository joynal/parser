package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"log"
	"os"

	"lambda-push-go/core"

	"github.com/mongodb/mongo-go-driver/bson"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DB_NAME")
	stream := os.Getenv("PARSER_TOPIC")

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

	notificationID, _ := primitive.ObjectIDFromHex("5cb854c0b565fc06edf6ee64")
	err = db.Collection("notifications").FindOne(ctx, bson.M{"_id": notificationID}).Decode(&notification)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Collection("notificationaccounts").FindOne(ctx, bson.D{}).Decode(&notificationAccount)

	if err != nil {
		log.Fatal(err)
	}

	processed, _ := json.Marshal(core.ProcessedNotification{
		ID:            notification.ID,
		SiteID:        notification.SiteID,
		TimeToLive:    notification.TimeToLive,
		LaunchURL:     notification.LaunchURL,
		Message:       notification.Messages[0],
		Browser:       notification.Browsers,
		HideRules:     notification.HideRules,
		TotalSent:     notification.TotalSent,
		SendTo:        notification.SendTo,
		IsAtLocalTime: false,

		// notification account data
		VapidDetails: notificationAccount.VapidDetails,
	})

	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}

	topic := client.Topic(stream)

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(processed),
	})

	id, err := result.Get(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Sent msg ID: %v\n", id)
}

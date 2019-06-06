package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
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

	// find subscribers and
	subscriberCol := db.Collection("notificationsubscribers")
	siteId, _ := primitive.ObjectIDFromHex("5c82424627ff1506951b7fbb")
	cur, err := subscriberCol.Find(ctx, bson.M{
		"siteId": siteId,
		"status": "subscribed",
	})

	if err != nil {
		log.Fatal(err)
	}
	// Close the cursor once finished
	defer cur.Close(ctx)

	// Iterate through the cursor
	for cur.Next(ctx) {
		var sub core.Subscriber
		err := cur.Decode(&sub)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("subscriberId: ", sub.ID)
	}
}

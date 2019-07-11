package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"io/ioutil"
	"log"
	"os"
	"parser/core"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/bson"
)

func main() {
	start := time.Now()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// prepare configs
	dbUrl := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DB_NAME")
	senderTopic := os.Getenv("SENDER_TOPIC")

	// Db connection stuff
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx = context.WithValue(ctx, core.DbURL, dbUrl)
	db, err := core.ConfigDB(ctx, dbName)
	if err != nil {
		log.Fatalf("database configuration failed: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	// process notification json object to struct
	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}

	// topic configs
	topic := client.Topic(senderTopic)

	// wait group for finishing all goroutines
	var wg sync.WaitGroup
	var notification core.ProcessedNotification
	largeJsonObjectStr, _ := ioutil.ReadFile("./test/message.json")

	err = json.Unmarshal(largeJsonObjectStr, &notification)
	if err != nil {
		fmt.Println("json error: ", err )
	}

	webPushOptions := webpush.Options{
		Subscriber:      "https://joynal.com/",
		VAPIDPublicKey:  notification.VapidDetails.VapidPublicKeys,
		VAPIDPrivateKey: notification.VapidDetails.VapidPrivateKeys,
		TTL:             notification.TimeToLive,
	}

	notificationPayload := core.NotificationPayload{
		ID:        notification.ID,
		LaunchURL: notification.LaunchURL,
		Message:   notification.Message,
		Browser:   notification.Browser,
		HideRules: notification.HideRules,
		Actions:   notification.Actions,
	}

	notificationPayloadStr, _ := json.Marshal(notificationPayload)
	subscriberCol := db.Collection("notificationsubscribers")

	// get subscriber
	id, _ := primitive.ObjectIDFromHex("5ce6c3bafee270615db8352c")
	var subscriber core.Subscriber
	err = subscriberCol.FindOne(ctx, bson.M{"_id": id}).Decode(&subscriber)

	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go sendDataToTopic(core.SubscriberPayload{
		PushEndpoint: subscriber.PushEndpoint,
		Data:         string(notificationPayloadStr),
		Options:      webPushOptions,
		SubscriberID: subscriber.ID,
	}, ctx, topic, &wg)

	wg.Wait()

	// stop all topic's go routines
	topic.Stop()

	fmt.Println("elapsed:", time.Since(start))
}

func sendDataToTopic(subscriber core.SubscriberPayload, ctx context.Context, topic *pubsub.Topic, wg *sync.WaitGroup) {
	defer wg.Done()

	jsonData, _ := json.Marshal(subscriber)

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(jsonData),
	})
	id, err := result.Get(ctx)

	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Printf("Sent msg ID: %v\n", id)
}

package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"os"
	"parser/core"
	"sync"
	"time"
)

func main() {
	err := godotenv.Load()
		if err != nil {
		log.Fatal("Error loading .env file")
	}

	// prepare configs
	dbUrl := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DB_NAME")
	parserTopic := os.Getenv("PARSER_TOPIC")
	senderTopic := os.Getenv("SENDER_TOPIC")
	arraySize := 1000

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

	var notification core.ProcessedNotification
	sub := client.Subscription(parserTopic)
	cctx, _ := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		err = json.Unmarshal(msg.Data, &notification)
		if err != nil {
			fmt.Println(err)
		}

		// Lets prepare subscriber query
		query := bson.M{
			"siteId": notification.SiteID,
			"status": "subscribed",
		}

		// if notification have timezone
		if notification.Timezone != "" {
			query["timezone"] = notification.Timezone
		}

		webPushOptions := webpush.Options{
			Subscriber:      "https://omnikick.com/",
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

		// find subscribers and
		cur, err := subscriberCol.Find(ctx, query)
		if err != nil {
			log.Fatal(err)
		}
		// Close the cursor once finished
		defer cur.Close(ctx)

		// topic configs
		topic := client.Topic(senderTopic)
		topic.PublishSettings.CountThreshold = 1
		topic.PublishSettings.DelayThreshold = 3 * time.Second

		// wait group for finishing all goroutines
		var wg sync.WaitGroup
		start := time.Now()
		counter := 0

		// Iterate through the cursor
		var subscribers []core.SubscriberPayload
		for cur.Next(ctx) {
			var elem core.Subscriber
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatalln("encode err:", err)
			}

			subscribers = append(subscribers, core.SubscriberPayload{
				PushEndpoint: elem.PushEndpoint,
				Data:         string(notificationPayloadStr),
				Options:      webPushOptions,
				SubscriberID: elem.ID,
			})

			if len(subscribers) == arraySize {
				counter++
				fmt.Println("called from inside: ", counter)
				wg.Add(1)
				go sendDataToTopic(subscribers, ctx, topic, &wg)
				notification.TotalSent += arraySize
				subscribers = nil
			}
		}

		if len(subscribers) > 0 {
			fmt.Println("called from outside: ", counter)
			wg.Add(1)
			go sendDataToTopic(subscribers, ctx, topic, &wg)
			notification.TotalSent += len(subscribers)
			subscribers = nil
		}

		wg.Wait()
		fmt.Println("elapsed:", time.Since(start))

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		// stop all topic's go routines
		topic.Stop()

		// update notification stats
		updateQuery := bson.M{"updatedAt": time.Now()}
		updateQuery["totalSent"] = notification.TotalSent

		if notification.IsAtLocalTime == false {
			updateQuery["isProcessed"] = "done"
		}

		notificationCol := db.Collection("notifications")
		_, _ = notificationCol.UpdateOne(ctx, bson.M{"_id": notification.ID}, bson.M{"$set": updateQuery})
	})

	if err != nil {
		fmt.Println(err)
	}
}

func sendDataToTopic(subscribers []core.SubscriberPayload, ctx context.Context, topic *pubsub.Topic, wg *sync.WaitGroup) {
	defer wg.Done()

	jsonData, _ := json.Marshal(subscribers)

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(jsonData),
	})
	id, err := result.Get(ctx)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Sent msg ID: %v\n", id)
}

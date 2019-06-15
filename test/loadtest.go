package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()

	// process notification json object to struct
	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}

	// topic configs
	senderTopic := "notification"
	topic := client.Topic(senderTopic)
	defer topic.Stop()

	// topic settings
	topic.PublishSettings.DelayThreshold = 1 * time.Millisecond
	topic.PublishSettings.CountThreshold = 1000
	topic.PublishSettings.ByteThreshold = 2e6

	// wait group for finishing all goroutines
	var wg sync.WaitGroup
	defer wg.Wait()

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go sendDataToTopic(largeJsonObjectStr, ctx, topic, &wg)
		// add throttling, so buffered limit do not cross
		time.Sleep(1 * time.Millisecond)
	}
}

func sendDataToTopic(notification string, ctx context.Context, topic *pubsub.Topic, wg *sync.WaitGroup) {
	defer wg.Done()

	jsonData, _ := json.Marshal(notification)

	fmt.Println("number of bytes: ", len(jsonData))

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(jsonData),
	})
	id, err := result.Get(ctx)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Sent msg ID: %v\n", id)
}

var largeJsonObjectStr = `{
    "_id" : "5c82455744cd0f069b35daa6",
    "priority" : "high",
    "timeToLive" : 259200,
    "totalSent" : 55,
    "totalDeliver" : 0,
    "totalShow" : 0,
    "totalError" : 0,
    "totalClick" : 0,
    "totalClose" : 0,
    "isAtLocalTime" : false,
    "isProcessed" : "done",
    "isSchedule" : true,
    "timezonesCompleted" : [],
    "isDeleted" : false,
    "receivers" : [],
    "actions" : [],
    "fromRSSFeed" : false,
    "siteId" : "5c82424627ff1506951b7fbb",
    "messages" : [
        {
            "_id" : "5c82455744cd0f069b35daa7",
            "title" : "Load test, pls ignore - 1",
            "message" : "Load test, pls ignore - 1",
            "language" : "en"
        }
    ],
    "browsers" : [
        {
            "iconUrl" : "https://cdn.testsite.com/assets/img/logo.png",
            "imageUrl" : "",
            "badge" : "",
            "vibration" : false,
            "isActive" : true,
            "isEnabledCTAButton" : false,
            "browserName" : "chrome"
        }
    ],
    "launchUrl" : "https://joynal.github.io",
    "userId" : "5c82424427ff1506951b7fb8",
    "sentAt" : "2019-03-29T16:34:00.040Z"
}`

package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"log"
	"os"
	"parser/core"
	"sync"
	"time"
)

func main() {
	// Db connection stuff
	ctx := context.Background()

	// process notification json object to struct
	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}

	// topic configs
	senderTopic := "notification"
	topic := client.Topic(senderTopic)
	topic.PublishSettings.DelayThreshold = 1 * time.Second
	topic.PublishSettings.CountThreshold = 1
	topic.PublishSettings.ByteThreshold = 3e6

	// stop all topic's go routines
	defer topic.Stop()

	// wait group for finishing all goroutines
	var wg sync.WaitGroup
	defer wg.Wait()

	var subscribers []core.SubscriberPayload
	subscriberId, _ := primitive.ObjectIDFromHex("5c9dea2e6ffa8676f90163d6")
	for i := 0; i < 1000000; i++ {
		subscribers = append(subscribers, core.SubscriberPayload{
			PushEndpoint: `{"endpoint":"https://fcm.googleapis.com/fcm/send/C5GvDOw9Nnf:APA91bESNu5qsIA484DSFWyuDLEgMHdAJf45IwMua9lknXrhAzQCrLcN-ZWfT8GE-_kxNR6MiCq1tfPr1aKWH8bVFNm6bmtDY-xHug-B76h6IqwemtB9tnlPsTqlr9A8ZcvA3dZzlxMc","expirationTime":null,"keys":{"p256dh":"BHk1DzprVgT26pIBTc3gsm-xE1m-DZzZcn_xAnvEpGKBMkja3V5rQsFQuQ7wlJV6I0A2P5LVHtjhp7lYZPsoQ8E","auth":"mrLLfPc_dIlwsO521ix1bQ"}}`,
			Data: `{
    "_id" : ObjectId("5c82455744cd0f069b35daa6"),
    "sendTo" : {
        "allSubscriber" : true,
        "segments" : []
    },
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
    "siteId" : ObjectId("5c82424627ff1506951b7fbb"),
    "messages" : [ 
        {
            "_id" : ObjectId("5c82455744cd0f069b35daa7"),
            "title" : "Load test, pls ignore - 1",
            "message" : "Load test, pls ignore - 1",
            "language" : "en"
        }
    ],
    "browsers" : [ 
        {
            "iconUrl" : "https://cdn.omnikick.com/assets/img/Logo_Smile_blue.png",
            "imageUrl" : "",
            "badge" : "",
            "vibration" : false,
            "isActive" : true,
            "isEnabledCTAButton" : false,
            "_id" : ObjectId("5c82455744cd0f069b35daa9"),
            "browserName" : "chrome",
            "actions" : [ 
                {
                    "_id" : ObjectId("5c82455744cd0f069b35daaa"),
                    "title" : "",
                    "url" : "",
                    "action" : "button1"
                }
            ]
        }, 
        {
            "iconUrl" : "https://cdn.omnikick.com/assets/img/Logo_Smile_blue.png",
            "imageUrl" : "",
            "badge" : "",
            "vibration" : false,
            "isActive" : true,
            "isEnabledCTAButton" : false,
            "_id" : ObjectId("5c82455744cd0f069b35daa8"),
            "browserName" : "firefox",
            "actions" : []
        }
    ],
    "launchUrl" : "https://joynal.github.io",
    "userId" : ObjectId("5c82424427ff1506951b7fb8"),
    "sentAt" : ISODate("2019-03-29T16:34:00.040Z")
}`,
			Options: webpush.Options{
				Subscriber:      "https://omnikick.com/",
				VAPIDPublicKey:  "BAPJRulAlnikXjrq7lEMKuhYucDIDhSDmGGtaUOBb31zg_PQZxr-suQk4LIeseAyROWwVMH4AONjK8yGb6cz5fs",
				VAPIDPrivateKey: "MLhfR_zQ8jg2D1eMy-zNmFLJ3QotNKlYfhbkmlox5vw",
				TTL:             259200,
			},
			SubscriberID: subscriberId,
		})

		if len(subscribers) == 1000 {
			wg.Add(1)
			go sendDataToTopic(subscribers, ctx, topic, &wg)
			subscribers = nil
		}
	}
}

func sendDataToTopic(subscribers []core.SubscriberPayload, ctx context.Context, topic *pubsub.Topic, wg *sync.WaitGroup) {
	defer wg.Done()

	jsonData, _ := json.Marshal(subscribers)

	fmt.Println("number of bytes: ", len(jsonData))

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(jsonData),
	})
	id, err := result.Get(ctx)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Sent msg ID: %v\n", id)
	// add throttling, so buffer limit do not exist
	time.Sleep(1 * time.Second)
}

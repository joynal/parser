package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

	topic.PublishSettings.CountThreshold = 1000
	topic.PublishSettings.BufferedByteLimit = 2e9

	// wait group for finishing all goroutines
	var wg sync.WaitGroup
	// defer wg.Wait()
	start := time.Now()

	largeJsonObjectStr, _ := ioutil.ReadFile("./test/message.json")

	for i := 0; i < 1000000; i++ {
		result := topic.Publish(ctx, &pubsub.Message{
			Data: largeJsonObjectStr,
		})

		wg.Add(1)
		go getResult(result, ctx, &wg)

		// add throttling, so buffered limit do not cross
		// time.Sleep(1 * time.Millisecond)
	}

	wg.Wait()

	fmt.Println("elapsed:", time.Since(start))
}

func getResult(result *pubsub.PublishResult, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	id, err := result.Get(ctx)

	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Printf("Sent msg ID: %v\n", id)
}

package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, mustGetenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}

	topicName := "notification"
	topic := client.Topic(topicName)

	for i:=0; i<10 ; i++ {
		result := topic.Publish(ctx, &pubsub.Message{
			Data: []byte(fmt.Sprintf("sample test - %d", i)),
		})

		id, err := result.Get(ctx)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Sent msg ID: %v\n", id)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

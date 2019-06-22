package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "adept-mountain-238503")
	if err != nil {
		log.Fatal(err)
	}

	topicName := "raw-notification"
	topic := client.Topic(topicName)

	message, err := ioutil.ReadFile("./test/message.json")
	if err != nil {
		fmt.Println(err)
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: message,
	})

	id, err := result.Get(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Sent msg ID: %v\n", id)
}

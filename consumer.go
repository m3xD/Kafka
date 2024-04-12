package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"time"
)

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	if err != nil {
		log.Fatalln("Something went wrong", err)
	}

	c.SubscribeTopics([]string{"demo_topic"}, nil)

	run := true

	for run {
		msg, err := c.ReadMessage(time.Second)

		if err == nil {
			fmt.Print("Message is", msg)
		} else if !err.(kafka.Error).IsFatal() {
			fmt.Print("Consumer error")
		}
	}

	c.Close()
}

package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Define Kafka consumer configuration
	config := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	}

	// Create a new Kafka consumer
	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		fmt.Errorf("Error creating consumer: %v", err)
	}

	// Subscribe to the topic
	err = consumer.Subscribe("clickview", nil)
	if err != nil {
		fmt.Errorf("Error subscribing: %v", err)
	}

	// Consume messages
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	// Close the consumer when done
	consumer.Close()
}

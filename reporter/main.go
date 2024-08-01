package main

import (
	"common"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	DB                 *gorm.DB
	msgChan            = make(chan common.EventServiceApiModel, 100)
	batchMapImpression = make(map[int]int)
	batchMapClick      = make(map[int]int)
	mu                 sync.Mutex
)

func main() {
	if err := OpenDbConnection(); err != nil {
		fmt.Printf("open db connection failed, err:%v\n", err)
	}
	if err := DB.AutoMigrate(&aggrClick{}); err != nil {
		fmt.Printf("auto migrate failed, err:%v\n", err)
	}
	if err := DB.AutoMigrate(&aggrImpression{}); err != nil {
		fmt.Printf("auto migrate failed, err:%v\n", err)
	}
	// Define Kafka consumer configuration
	go ComsumeMessageKafka()
	go handlebatch(msgChan)

	select {}
}

func handlebatch(ch chan common.EventServiceApiModel) {
	for {
		select {
		case event := <-ch:
			mu.Lock()
			if event.IsClicked == true {
				batchMapClick[event.AdId]++
			} else {
				batchMapImpression[event.AdId]++
			}
			mu.Unlock()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func ComsumeMessageKafka() {
	config := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	}

	// Create a new Kafka consumer
	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		fmt.Println("Error creating consumer: %v", err)
	}

	// Subscribe to the topic
	err = consumer.Subscribe("clickview", nil)
	if err != nil {
		fmt.Println("Error subscribing: %v", err)
	}

	// Consume messages
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			var infoImpressionClick common.EventServiceApiModel
			if err := json.Unmarshal(msg.Value, &infoImpressionClick); err != nil {
				fmt.Printf("Error unmarshalling message: %v", err)
				continue
			}
			msgChan <- infoImpressionClick
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	// Close the consumer when done
	consumer.Close()
}

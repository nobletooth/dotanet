package main

import (
	"flag"
	_ "github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	KafkaTopic  = flag.String("kafka-topic", "my-topic", "Kafka topic")
	kafkaBroker = flag.String("kafka-broker", "localhost:9092", "Kafka broker")
)

func main() {
	flag.Parse()

}

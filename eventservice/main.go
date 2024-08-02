package main

import (
	"common"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	Db               *gorm.DB
	producer         *kafka.Producer
	ch               = make(chan common.EventServiceApiModel, 10)
	impressionEvents []common.EventServiceApiModel
	eventsMutex      sync.Mutex
)

func main() {
	flag.Parse()
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": *kafkaendpoint})
	if err != nil {
		fmt.Printf("\nerror opening kafka connection: %v\n", err)
	}
	go panelApiCall(ch, producer)
	if db, err := OpenDbConnection(); err != nil {
		fmt.Printf("\nerror opening db connection: %v\n", err)
	} else {
		Db = db
	}
	go cleanOldEvents()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // Change to your frontend domain
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/click/:adv/:pub/:clickid/:impressionid", clickHandler())
	router.GET("/impression/:adv/:pub/:impressionid", impressionHandler())

	router.Run(*EventserviceUrl)

}

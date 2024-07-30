package main

import (
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
	"gorm.io/gorm"
	"time"
)

var Db *gorm.DB
var ch = make(chan common.EventServiceApiModel, 10)
var p *kafka.Producer

func main() {
	go panelApiCall(ch)
	_, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		fmt.Errorf("error opening kafka connection: %v", err)

	}
	flag.Parse()
	if db, err := OpenDbConnection(); err != nil {
		fmt.Errorf("error opening db connection: %v", err)
	} else {
		Db = db
	}
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // Change to your frontend domain
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/click/:adv/:pub", clickHandler())
	router.GET("/impression/:adv/:pub", impressionHandler())

	router.Run(*EventserviceUrl)

}

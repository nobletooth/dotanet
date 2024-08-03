package main

import (
	"common"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
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
	impressionEvents []common.UrlImpressionParameters
	eventsMutex      sync.Mutex
	userClickTracker = make(map[string][]time.Time)
	userClickMutex   sync.Mutex
)

func main() {
	flag.Parse()
	go cleanOldEvents()
	go cleanOldClickData()
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
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/click/:encryptedClickParams", clickHandler())
	router.GET("/impression/:encryptedImpressionParams", impressionHandler())

	router.Run(*EventserviceUrl)

}

func decrypt(encodedData string) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher([]byte(*secretKey))
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func cleanOldClickData() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			userClickMutex.Lock()
			for userID, clicks := range userClickTracker {
				cutoff := time.Now().Add(*userClickCutoff)
				userClickTracker[userID] = filterRecentClicks(clicks, cutoff)
				if len(userClickTracker[userID]) == 0 {
					delete(userClickTracker, userID)
				}
			}
			userClickMutex.Unlock()
		}
	}
}

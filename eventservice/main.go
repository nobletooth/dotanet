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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	Db               *gorm.DB
	ch               = make(chan common.EventServiceApiModel, 10)
	impressionEvents []common.EventServiceApiModel
	eventsMutex      sync.Mutex
)

func main() {
	go panelApiCall(ch)
	flag.Parse()
	if db, err := OpenDbConnection(); err != nil {
		fmt.Errorf("error opening db connection: %v", err)
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

package main

import (
	"common"
	"flag"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	allAds                            []common.AdWithMetrics
	AdserverUrl                       = flag.String("adserverurl", ":8081", "ad server port")
	EventServiceUrl                   = flag.String("eventserviceurl", "http://localhost:8082", "ad server port")
	PanelUrlGetAllAds                 = flag.String("panelurlads", "http://localhost:8085", "panel url")
	PanelUrlPicUrl                    = flag.String("panelurlpic", "http://localhost:8085", "panel url")
	NewAdImpressionThreshold          = flag.Int64("newAdTreshold", 5, "Impression threshold for considering an ad as new")
	NewAdSelectionProbability         = flag.Float64("newAdProb", 0.25, "Probability of selecting a new ad")
	ExperiencedAdSelectionProbability = flag.Float64("expAdProb", 0.75, "Probability of selecting a exprienced ad")
	secretKey                         = flag.String("secretkey", "X9K3jM5nR7pL2qT8vW1cY6bF4hG0xA9E", "secret key")
)
var config = cors.Config{
	AllowAllOrigins:  true,
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	AllowHeaders:     []string{"*"},
	ExposeHeaders:    []string{"*"},
	AllowCredentials: false,
	MaxAge:           12 * time.Hour,
}

func main() {
	flag.Parse()
	router := gin.Default()
	router.Use(cors.New(config))

	router.GET("/getadinfo/:pubID", GetAdHandler)

	go GetAdsListPeriodically()
	fmt.Println("Server running on port" + *AdserverUrl)
	router.Run(*AdserverUrl)
}

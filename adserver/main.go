package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
)

var (
	allAds                            []common.AdWithMetrics
	AdserverPort                      = flag.String("adserverport", ":8081", "ad server port")
	PanelUrl                          = flag.String("panelurl", "http://localhost:8080", "panel url")
	NewAdImpressionThreshold          = flag.Int64("newAdTreshold", 5, "Impression threshold for considering an ad as new")
	NewAdSelectionProbability         = flag.Float64("newAdProb", 0.25, "Probability of selecting a new ad")
	ExperiencedAdSelectionProbability = flag.Float64("expAdProb", 0.75, "Probability of selecting a exprienced ad")
)

func main() {
	flag.Parse()
	router := gin.Default()
	router.GET("/getad/:pubID", GetAdsHandler)
	router.GET("/getadinfo/:pubID", GetAdHandler)
	go GetAdsListPeriodically()
	fmt.Println("Server running on port" + *AdserverPort)
	router.Run(*AdserverPort)
}

func GetAdsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, allAds)
}

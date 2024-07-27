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
	AdserverUrl                       = flag.String("adserverurl", "localhost:8080", "ad server url")
	PanelUrl                          = flag.String("panelurl", "localhost:8081", "panel url")
	NewAdImpressionThreshold          = flag.Int64("newAdTreshold", 5, "Impression threshold for considering an ad as new")
	NewAdSelectionProbability         = flag.Float64("newAdProb", 0.25, "Probability of selecting a new ad")
	ExperiencedAdSelectionProbability = flag.Float64("expAdProb", 0.75, "Probability of selecting a exprienced ad")
)

func main() {
	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/getad/:pubID", GetAdsHandler)
	go GetAdsListPeriodically()
	fmt.Println("Server running on port" + *AdserverUrl)
	router.Run(*AdserverUrl)
}

func GetAdsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, allAds)
}

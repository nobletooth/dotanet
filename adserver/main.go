package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
)

var allAds []common.AdWithMetrics
var PanelPort string
var AdserverPort string

func init() {
	flag.StringVar(&AdserverPort, "adserverport", "8080", "ad server port")
	flag.StringVar(&PanelPort, "panelport", "8081", "panel port")
}

func main() {
	router := gin.Default()
	router.GET("/getad/:pubID", GetAdsHandler)
	go GetAdsListPeriodically()
	fmt.Println("Server running on port 8080")
	router.Run(":" + AdserverPort)
}

func GetAdsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, allAds)
}

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
)

var (
	allAds       []common.AdWithMetrics
	AdserverPort = flag.String("adserverport", "8080:", "ad server port")
	PanelUrl     = flag.String("panelurl", "http://localhost:8081", "panel url")
)

func main() {
	flag.Parse()
	router := gin.Default()
	router.GET("/getad/:pubID", GetAdsHandler)
	go GetAdsListPeriodically()
	fmt.Println("Server running on port" + *AdserverPort)
	router.Run(*AdserverPort)
}

func GetAdsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, allAds)
}

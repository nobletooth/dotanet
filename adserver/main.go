package main

import (
	"fmt"
	"net/http"

	"example.com/dotanet/common"
	"github.com/gin-gonic/gin"
)

var allAds []common.AdInfo

func main() {
	router := gin.Default()
	router.GET("/ads/", GetAdsHandler)
	go GetAdsListPeriodically(allAds)
	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}

func GetAdsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, allAds)
}

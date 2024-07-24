package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
)

var allAds []common.AdInfo

func main() {
	router := gin.Default()
	router.GET("/getad/:pubID", GetAdsHandler)
	go GetAdsListPeriodically()
	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}

func GetAdsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, allAds)
}

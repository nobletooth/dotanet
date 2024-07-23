package main

import (
	"fmt"

	"example.com/dotanet/common"
	"github.com/gin-gonic/gin"
)

var allAds []common.AdInfo

func main() {
	router := gin.Default()
	go GetAdsListPeriodically(allAds)
	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}

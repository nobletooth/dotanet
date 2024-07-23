package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/tree/main/common"
)

var allAds []common.AdInfo

func main() {
	router := gin.Default()
	go GetAdsListPeriodically()
	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}

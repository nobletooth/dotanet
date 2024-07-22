package main

import (
	"example.com/dotanet/panel/advertiser"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	if err := advertiser.NewDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	router := gin.Default()
	router.HTMLRender = advertiser.LoadTemplates("./templates")
	router.GET("/", advertiser.ListAdvertisers)
	router.GET("/advertisers/new", advertiser.NewAdvertiserForm)
	router.POST("/advertisers", advertiser.CreateAdvertiser)
	router.GET("/advertisers/:id", advertiser.GetAdvertiserCredit)
	router.GET("/advertisers/:id/ads", advertiser.ListAdsByAdvertiserHandler)
	router.GET("/ads/new", advertiser.CreateAdForm)
	router.POST("/ads", advertiser.CreateAdHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

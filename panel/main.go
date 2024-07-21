package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"example.com/dotanet/panel/advertiser"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	flag.Parse()

	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Close()

	if err := database.AutoMigrate(&advertiser.AdvertiserEntity{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	advertiserService := advertiser.NewAdvertiserService(database)

	router := gin.Default()
	setupRouter(router)

	router.POST("/advertisers", createAdvertiserHandler(advertiserService))
	router.GET("/advertisers", listAdvertisersHandler(advertiserService))
	router.GET("/advertisers/:id", getAdvertiserCreditHandler(advertiserService))

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func createAdvertiserHandler(service advertiser.AdvertiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto advertiser.AdvertiserRequestDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		service.CreateAdvertiserEntity(dto)
		c.JSON(http.StatusOK, gin.H{"craeted": "ok"})
	}
}

func listAdvertisersHandler(service advertiser.AdvertiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		advertisers := service.ListAllAdvertiserEntity()
		c.JSON(http.StatusOK, advertisers)
	}
}

func getAdvertiserCreditHandler(service advertiser.AdvertiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid advertiser ID"})
			return
		}

		credit, err := service.GetCreaditOfAdvertiser(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"credit": credit})
	}
}

func setupRouter(db *gorm.DB, router *gin.Default) {
	router.GET("/ads/new", advertiser.RenderCreateAdForm)
}

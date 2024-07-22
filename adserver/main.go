package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Ad struct {
	Id           uint    `gorm:"column:id;primary_key"`
	Title        string  `gorm:"column:title"`
	Image        string  `gorm:"column:image"`
	Price        float64 `gorm:"column:price"`
	Status       bool    `gorm:"column:status"`
	Clicks       int     `gorm:"column:clicks"`
	Impressions  int     `gorm:"column:impressions"`
	Url          string  `gorm:"column:url"`
	AdvertiserId uint64  `gorm:"foreignKey:AdvertiserId"`
}

var allAds []Ad

func main() {
	router := gin.Default()
	var db *gorm.DB
	r := multitemplate.NewRenderer()
	r.AddFromFiles("ads_list", "./templates/ads_list.html")
	fmt.Println("Server running on port 8080")
	go GetAdsListPeriodically(db, router)
	router.Run(":8080")
}

func GetAdsListPeriodically(db *gorm.DB, router *gin.Engine) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			router.GET("/ads/", func(c *gin.Context) {
				c.JSON(http.StatusOK, allAds)
				c.HTML(http.StatusOK, "ads_list.html", gin.H{
					"title": "All Ads",
					"ads":   allAds,
				})
			})

			for _, ad := range allAds {
				log.Printf("Ad ID: %d, Title: %s", ad.Id, ad.Title)
			}
			log.Println("Ads list fetched successfully")
		}
	}
}


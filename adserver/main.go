package main

import (
	"log"
	"net/http"
	"time"

	"example.com/dotanet/panel"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	router := gin.Default()
	setupRouter(router)
	go GetAdsListPeriodicly()
	router.Run(":8080")
}

func setupRouter(router *gin.Default) {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("ads_list", "./templates/ads_list.html")
	router.GET("/ads", AdsList)
}

func GetAdsListPeriodicly(db *gorm.DB, r *gin.Default) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c := &gin.Context{}
			response, err := r.HandleContext(c)
			if err != nil {
				log.Println("Error making GET request:", err)
				continue
			}
			log.Println("GET request successful, status:", response.Status())
		}
	}
}

func AdsList(c *gin.Context) {
	ads := panel.ListAllAds(c)
	c.HTML(http.StatusOK, "ads_list.html", gin.H{
		"title": "All Ads",
		"ads":   ads,
	})
}

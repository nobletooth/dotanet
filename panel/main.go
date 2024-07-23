package main

import (
	"log"
	"net/http"

	"example.com/dotanet/common"
	"example.com/dotanet/panel/advertiser"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func LoadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", templatesDir+"/index.html")
	r.AddFromFiles("create_advertiser", templatesDir+"/create_advertiser.html")
	r.AddFromFiles("advertiser_credit", templatesDir+"/advertiser_credit.html")
	r.AddFromFiles("create_ad", templatesDir+"/create_ad.html")
	r.AddFromFiles("advertiser_ads", templatesDir+"/advertiser_ads.html")
	r.AddFromFiles("edit_ad", templatesDir+"/edit_ad.html")
	return r
}

func main() {
	if err := advertiser.NewDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	router := gin.Default()
	router.HTMLRender = LoadTemplates("./templates")
	router.GET("/", advertiser.ListAdvertisers)
	router.GET("/advertisers/new", advertiser.NewAdvertiserForm)
	router.POST("/advertisers", advertiser.CreateAdvertiser)
	router.GET("/advertisers/:id", advertiser.GetAdvertiserCredit)
	router.GET("/advertisers/:id/ads", advertiser.ListAdsByAdvertiserHandler)
	router.GET("/ads/new", advertiser.CreateAdForm)
	router.POST("/ads", advertiser.CreateAdHandler)
	router.GET("/ads/edit/:id", advertiser.EditAdForm)
	router.POST("/ads/update/:id", advertiser.UpdateAdHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

package main

import (
	"example.com/dotanet/panel/advertiser"
	"example.com/dotanet/panel/common"
	"example.com/dotanet/panel/publisher"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func LoadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index.html", templatesDir+"/index.html")
	r.AddFromFiles("advertisers.html", templatesDir+"/advertisers.html")
	r.AddFromFiles("publishers.html", templatesDir+"/publishers.html")
	r.AddFromFiles("create_publisher.html", templatesDir+"/create_publisher.html")
	r.AddFromFiles("view.html", templatesDir+"/view.html")
	r.AddFromFiles("create_advertiser.html", templatesDir+"/create_advertiser.html")
	r.AddFromFiles("advertiser_credit.html", templatesDir+"/advertiser_credit.html")
	r.AddFromFiles("create_ad.html", templatesDir+"/create_ad.html")
	r.AddFromFiles("advertiser_ads.html", templatesDir+"/advertiser_ads.html")
	return r
}

func main() {
	if err := common.NewDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	router := gin.Default()
	router.HTMLRender = LoadTemplates("./templates")
	// Home route
	router.GET("/", homeHandler)

	// Advertiser routes
	router.GET("/advertisers", advertiser.ListAdvertisers)
	router.GET("/advertisers/new", advertiser.NewAdvertiserForm)
	router.POST("/advertisers", advertiser.CreateAdvertiser)
	router.GET("/advertisers/:id", advertiser.GetAdvertiserCredit)
	router.GET("/advertisers/:id/ads", advertiser.ListAdsByAdvertiserHandler)
	router.GET("/ads/new", advertiser.CreateAdForm)
	router.POST("/ads", advertiser.CreateAdHandler)

	// Publisher routes
	router.GET("/publishers", publisher.ListPublishers)
	router.GET("/publishers/new", publisher.NewPublisherForm)
	router.POST("/publishers", publisher.CreatePublisherHandler)
	router.GET("/publishers/:id", publisher.ViewPublisherHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

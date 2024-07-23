package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/tree/main/panel/advertiser"
	"github.com/nobletooth/dotanet/tree/main/panel/common"
	"github.com/nobletooth/dotanet/tree/main/panel/publisher"
)

func LoadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", templatesDir+"/index.html")
	r.AddFromFiles("create_advertiser", templatesDir+"/create_advertiser.html")
	r.AddFromFiles("advertisers", templatesDir+"/advertisers.html") // Add this line
	r.AddFromFiles("advertiser_credit", templatesDir+"/advertiser_credit.html")
	r.AddFromFiles("create_ad", templatesDir+"/create_ad.html")
	r.AddFromFiles("advertiser_ads", templatesDir+"/advertiser_ads.html")
	r.AddFromFiles("publishers", templatesDir+"/publishers.html")
	r.AddFromFiles("create_publisher", templatesDir+"/create_publisher.html")
	r.AddFromFiles("view_publisher", templatesDir+"/view.html")
	return r
}

func main() {
	if err := common.NewDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	err := common.DB.AutoMigrate(&publisher.Publisher{})
	err = common.DB.AutoMigrate(&advertiser.Ad{}, &advertiser.Ad{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)

	}

	router := gin.Default()
	router.Use(cors.Default())

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
	router.POST("/ads/update", advertiser.UpdateAdHandler)
	router.GET("/ads/:id/picture", advertiser.LoadAdPictureHandler)

	// Publisher routes
	router.GET("/publishers", publisher.ListPublishers)
	router.GET("/publishers/new", publisher.NewPublisherForm)
	router.POST("/publishers", publisher.CreatePublisherHandler)
	router.GET("/publishers/:id", publisher.ViewPublisherHandler)
	router.GET("/publishers/:id/script", publisher.GetPublisherScript)

	// Ad Server
	router.GET("/ads/list/", advertiser.ListAllAds)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}

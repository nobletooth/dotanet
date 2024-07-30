package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
	"github.com/nobletooth/dotanet/panel/advertiser"
	"github.com/nobletooth/dotanet/panel/database"
	"github.com/nobletooth/dotanet/panel/publisher"
)

var config = cors.Config{
	AllowAllOrigins:  true,
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	AllowHeaders:     []string{"*"},
	ExposeHeaders:    []string{"*"},
	AllowCredentials: false,
	MaxAge:           12 * time.Hour,
}

func eventservice(event common.EventServiceApiModel) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if event.IsClicked {
		clickedEvent := common.ClickedEvent{
			ID:           event.ClickID,
			ImpressionID: event.ImpressionID,
			Pid:          event.PubId,
			AdId:         event.AdId,
			Time:         event.Time,
		}
		result := tx.Create(&clickedEvent)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		ad, err := advertiser.FindAdById(clickedEvent.AdId)
		if err != nil {
			tx.Rollback()
			return err
		}
		if err := advertiser.HandleAdvertiserCredit(tx, ad); err != nil {
			tx.Rollback()
			return err
		}
		if err := publisher.HandlePublisherCreditWithTx(tx, ad, clickedEvent.Pid); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		viewedEvent := common.ViewedEvent{
			ID:   event.ImpressionID,
			Pid:  event.PubId,
			AdId: event.AdId,
			Time: event.Time,
		}

		result := tx.Create(&viewedEvent)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func eventServerHandler(c *gin.Context) {
	var event common.EventServiceApiModel
	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := eventservice(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event processed successfully", "event": event})
}

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
	r.AddFromFiles("reports", templatesDir+"/reports.html")
	r.AddFromFiles("ad_reports", templatesDir+"/ad_reports.html")
	return r
}

func main() {
	flag.Parse()
	if err := database.NewDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	err := database.DB.AutoMigrate(&publisher.Publisher{})
	err = database.DB.AutoMigrate(&advertiser.Entity{}, &advertiser.Ad{})
	err = database.DB.AutoMigrate(&common.ClickedEvent{}, &common.ViewedEvent{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)

	}

	router := gin.Default()
	router.Use(cors.New(config))
	router.HTMLRender = LoadTemplates("./templates")
	// Home route
	router.GET("/", homeHandler)

	// Advertiser routes
	router.GET("/advertisers", advertiser.ListAdvertisers)
	router.GET("/advertisers/new", advertiser.NewAdvertiserForm)
	router.POST("/advertisers", advertiser.CreateAdvertiser)
	router.GET("/advertisers/:id", advertiser.GetAdvertiserCredit)
	router.GET("/advertisers/:id/ads", advertiser.ListAdsByAdvertiserHandler)
	router.GET("/advertisers/:id/reports", advertiser.GetAdvertiserAdReports)

	router.GET("/ads/new", advertiser.CreateAdForm)
	router.POST("/ads", advertiser.CreateAdHandler)
	router.POST("/ads/update", advertiser.UpdateAdHandler)
	router.GET("/ads/:id/picture", advertiser.LoadAdPictureHandler) //endpoint

	// Publisher routes
	router.GET("/publishers", publisher.ListPublishers)
	router.GET("/publishers/new", publisher.NewPublisherForm)
	router.POST("/publishers", publisher.CreatePublisherHandler)
	router.GET("/publishers/:id", publisher.ViewPublisherHandler)
	router.GET("/publishers/:id/script", publisher.GetPublisherScript)
	router.POST("/eventservice", eventServerHandler) //end point
	router.GET("/publishers/:id/reports", publisher.GetPublisherReports)

	// Ad server routes
	router.GET("/ads/list/", advertiser.ListAllAds) //endpoint

	if err := router.Run(*(database.PanelUrl)); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}

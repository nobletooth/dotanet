package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	sitenames        = []string{"digikala", "digiland", "samsung", "torob", "varzesh3"}
	PublisherService = flag.String("publisherservice", ":8083", "publisher service")
	BaseURL          = flag.String("baseurl", "95.217.125.139:8085", "Base URL for HTML files")
)

func main() {
	flag.Parse()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // Change to your frontend domain
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.LoadHTMLGlob("./html/*")

	router.GET("/:sitename", siteHandler(*BaseURL))

	router.Run(*PublisherService)
}

func siteHandler(baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sitename := c.Param("sitename")
		var siteExist = false
		for _, value := range sitenames {
			if value == sitename {
				siteExist = true
				break
			}
		}
		if !siteExist {
			c.String(http.StatusBadRequest, "error : this site does not exist.")
			return
		}
		htmladdress := fmt.Sprintf("%v.html", sitename)
		c.HTML(http.StatusOK, htmladdress, gin.H{
			"BaseURL": baseURL,
		})
	}
}

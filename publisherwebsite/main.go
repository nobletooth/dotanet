package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	sitenames        = []string{"digikala", "digiland", "samsung", "torob", "varzesh3"}
	PublisherService = flag.String("publisherservice", ":8083", "publisher service")
)

func main() {
	flag.Parse()
	router := gin.Default()
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))

	router.LoadHTMLGlob("./html/*")

	router.GET("/:sitename", siteHandler())

	router.Run(*PublisherService)
}

func siteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		sitename := c.Param("sitename")
		var siteExist = false
		for _, value := range sitenames {
			if value == sitename {
				siteExist = true
			}
		}
		if !siteExist {
			c.String(http.StatusBadRequest, "error : this site does not exist.")
		}
		htmladdress := fmt.Sprintf("%v.html", sitename)
		c.HTML(http.StatusOK, htmladdress, gin.H{})
	}
}

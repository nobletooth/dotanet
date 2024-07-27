package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	sitenames        = []string{"digikala", "digiland", "samsung", "torob", "varzesh3"}
	PublisherService = flag.String("publisherservice", ":8083", "publisher service")
)

func main() {
	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
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

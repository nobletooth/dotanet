package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var sitenames = []string{"varzesh3", "digikala", "samsung", "torob", "digiland"}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("./publisherwebsite/html/*")

	// digikalaHandler
	router.GET("/digikala", func(c *gin.Context) {
		c.HTML(http.StatusOK, "digikala/index.tmpl", gin.H{
			"title": "Digikala",
		})
	})

	router.GET("/", defaultHandler)
	router.POST("/", defaultHandler)

	router.Run(":6060")
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

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func main() {
	pgConnection()
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")

	// varzesh3Handler
	router.GET("/varzesh3", func(c *gin.Context) {
		c.HTML(http.StatusOK, "varzesh3/index.tmpl", gin.H{
			"title": "Varzesh3",
		})
	})

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

//func varzesh3Handler(content string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.String(200, "Welcome to Varzesh3!")
//	}
//}
//
//func digikalaHandler(content string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.String(200, "Welcome to Digikala!")
//	}
//}

// defaultHandler
func defaultHandler(c *gin.Context) {
	c.String(200, "Welcome to API server!")
}

func pgConnection() {

}

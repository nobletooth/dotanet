package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//pgConnection()
	router := gin.Default()
	router.LoadHTMLGlob("./publisherwebsite/html/*")

	router.GET("/torob", torobHandler())
	router.GET("/samsung", samsungHandler())
	router.GET("/digiland", digilandHandler())
	router.GET("/varzesh3", varzesh3Handler())
	router.GET("/digikala.js", digikalaHandler())

	router.Run(":6060")
}

func torobHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "torob.html", gin.H{})
	}
}

func samsungHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "samsung.html", gin.H{})
	}
}

func digilandHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "digiland.html", gin.H{})
	}
}

func varzesh3Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "varzesh3.html", gin.H{})
	}
}

func digikalaHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "digikala.js.html", gin.H{})
	}
}

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//pgConnection()
	router := gin.Default()
	router.LoadHTMLGlob("./html/*")


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

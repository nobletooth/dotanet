package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	pgConnection()
	router := gin.Default()

	router.GET("/")
	router.POST("/")

	router.Run(":6060")
}

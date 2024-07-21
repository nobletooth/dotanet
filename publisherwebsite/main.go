package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func main() {
	//pgConnection()
	router := gin.Default()

	router.GET("/")
	router.POST("/")

	router.Run(":6060")
}

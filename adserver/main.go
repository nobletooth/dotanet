package main

import (
	"example.com/dotanet/panel"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	setupRouter(router)
}

func setupRouter(router *gin.Default) {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("ads_list", "./templates/ads_list.html")
	router.GET("/ads", panel.ListAllAds)
}

package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
	"gorm.io/gorm"
)

var Db *gorm.DB
var ch = make(chan common.EventServiceApiModel, 10)

var config = cors.Config{
	AllowAllOrigins:  true,
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	AllowHeaders:     []string{"*"},
	ExposeHeaders:    []string{"*"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
}

func main() {
	flag.Parse()
	if db, err := OpenDbConnection(); err != nil {
		fmt.Errorf("error opening db connection: %v", err)
	} else {
		Db = db
	}
	router := gin.Default()
	router.Use(cors.New(config))

	router.GET("/click/:adv/:pub", clickHandler())
	router.GET("/impression/:adv/:pub", impressionHandler())

	router.Run(*EventservicePort)

	go panelApiCall(ch)
}

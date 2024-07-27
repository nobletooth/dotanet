package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
	"gorm.io/gorm"
)

var Db *gorm.DB
var ch = make(chan common.EventServiceApiModel, 10)

func main() {
	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	if db, err := OpenDbConnection(); err != nil {
		fmt.Errorf("error opening db connection: %v", err)
	} else {
		Db = db
	}
	router := gin.Default()

	router.GET("/click/:adv/:pub", clickHandler())
	router.GET("/impression/:adv/:pub", impressionHandler())

	router.Run(*EventserviceUrl)

	go panelApiCall(ch)
}

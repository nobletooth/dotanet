package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
	"net/http"
	"time"
)

func impressionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var impressionTime = time.Now()
		adv := c.Param("adv") // should decrypt adv and pub.
		pub := c.Param("pub") // should decrypt adv and pub.
		var updateApi = common.EventServiceApiModel{Time: impressionTime, PubId: pub, AdId: adv, IsClicked: false}
		ch <- updateApi
		c.String(http.StatusOK, "its ok!")
	}
}

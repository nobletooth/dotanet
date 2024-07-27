package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
	"net/http"
	"strconv"
	"time"
)

func impressionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var impressionTime = time.Now()
		adv := c.Param("adv") // should decrypt adv and pub.
		pub := c.Param("pub") // should decrypt adv and pub.
		pubInt, err := strconv.Atoi(pub)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		adInt, err := strconv.Atoi(adv)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var updateApi = common.EventServiceApiModel{Time: impressionTime, PubId: pubInt, AdId: adInt, IsClicked: false}
		ch <- updateApi
		c.String(http.StatusOK, "its ok!")
	}
}

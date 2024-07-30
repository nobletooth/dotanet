package main

import (
	"common"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func impressionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		impressionId := uuid.MustParse(c.Param("impressionid"))
		time, err := time.Parse(time.RFC3339, c.Param("time"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time"})
			return
		}
		var updateApi = common.EventServiceApiModel{Time: time,
			PubId: pubInt, AdId: adInt, IsClicked: false, ImpressionID: impressionId}
		ch <- updateApi
		c.String(http.StatusOK, "its ok!")
	}
}

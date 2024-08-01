package main

import (
	"common"
	"fmt"
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
			return
		}
		adInt, err := strconv.Atoi(adv)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		impressionId := uuid.MustParse(c.Param("impressionid")) //primary key : uuid
		fmt.Printf("\n\n\nImpressionId: %x\n\n\n", impressionId)

		impressionTime := time.Now()
		var updateApi = common.EventServiceApiModel{
			Time:         impressionTime,
			PubId:        pubInt,
			AdId:         adInt,
			IsClicked:    false,
			ImpressionID: impressionId,
		}
		eventsMutex.Lock()
		impressionEvents = append(impressionEvents, updateApi)
		eventsMutex.Unlock()
		ch <- updateApi
		c.String(http.StatusOK, "its ok!")
	}
}

func cleanOldEvents() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			currentTime := time.Now()
			eventsMutex.Lock()
			i := 0
			for _, event := range impressionEvents {
				if currentTime.Sub(event.Time) <= 30*time.Second {
					impressionEvents[i] = event
					i++
				}
			}
			impressionEvents = impressionEvents[:i]
			eventsMutex.Unlock()
		}
	}
}

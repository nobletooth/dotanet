package main

import (
	"common"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func impressionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// decryption
		encryptedImpressionParams := c.Param("encryptedImpressionParams")
		decryptedImpressionParams, _ := decrypt(encryptedImpressionParams)
		var impressionParams common.ViewedEvent
		json.Unmarshal(decryptedImpressionParams, &impressionParams)
		adID := impressionParams.AdId
		pubID := impressionParams.Pid
		impressionId := impressionParams.ID

		// deduplicate impression
		if !checkDuplicateImpression(impressionId) {

			impressionTime := time.Now()
			var updateApi = common.EventServiceApiModel{
				Time:         impressionTime,
				PubId:        adID,
				AdId:         pubID,
				IsClicked:    false,
				ImpressionID: impressionId,
			}

			eventsMutex.Lock()
			impressionEvents = append(impressionEvents, updateApi)
			eventsMutex.Unlock()
			ch <- updateApi
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Duplicate impression"})
			return
		}
		c.String(http.StatusOK, "its ok!")
	}
}

func checkDuplicateImpression(impressionId uuid.UUID) bool {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	for _, event := range impressionEvents {
		if event.ImpressionID == impressionId && event.IsClicked {
			return true
		}
	}
	return false
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

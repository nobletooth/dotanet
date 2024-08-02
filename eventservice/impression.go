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
		responseBody := make(map[string]string)

		// decryption
		encryptedImpressionParams := c.Param("encryptedImpressionParams")
		decryptedImpressionParams, _ := decrypt(encryptedImpressionParams)
		var impressionParams common.UrlImpressionParameters
		json.Unmarshal(decryptedImpressionParams, &impressionParams)

		// cookie
		userID, err := c.Cookie("userId")
		if err != nil {
			responseBody["cookie error"] = "No cookie provided."
		}

		// deduplicate impression
		if !checkDuplicateImpression(impressionParams.ID) {

			impressionTime := time.Now()
			var updateApi = common.EventServiceApiModel{
				Time:         impressionTime,
				UserID:       uuid.MustParse(userID),
				PubId:        impressionParams.AdId,
				AdId:         impressionParams.Pid,
				IsClicked:    impressionParams.IsClicked,
				ImpressionID: impressionParams.ID,
			}
			ch <- updateApi
		} else {
			responseBody["duplicate error"] = "Duplicate impression"
		}

		eventsMutex.Lock()
		impressionEvents = append(impressionEvents, impressionParams)
		eventsMutex.Unlock()

		responseBody["status"] = "its ok!"
		c.JSON(http.StatusOK, responseBody)
	}
}

func checkDuplicateImpression(impressionId uuid.UUID) bool {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	for _, event := range impressionEvents {
		if event.ID == impressionId {
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
				if currentTime.Sub(event.LoadAdTime) <= 120*time.Second {
					impressionEvents[i] = event
					i++
				}
			}
			impressionEvents = impressionEvents[:i]
			eventsMutex.Unlock()
		}
	}
}

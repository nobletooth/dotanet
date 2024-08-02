package main

import (
	"bytes"
	"common"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func clickHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		responseBody := make(map[string]string)

		// decryption
		encryptedClickParams := c.Param("encryptedClickParams")
		decryptedClickParams, _ := decrypt(encryptedClickParams)
		var clickParams common.UrlClickParameters
		json.Unmarshal(decryptedClickParams, &clickParams)

		// refresh and click tracker
		userID, err := c.Cookie("userId")
		if err != nil {
			responseBody["cookie error"] = "No cookie provided."
		}
		userClickMutex.Lock()
		defer userClickMutex.Unlock()

		currentTime := time.Now()

		// cutoff : 5 minutes ago
		cutoff := currentTime.Add(*userClickCutoff)
		userClickTracker[userID] = filterRecentClicks(userClickTracker[userID], cutoff)

		if len(userClickTracker[userID]) >= *limitUserClick {
			responseBody["cookie error"] = "Same cookie clicked more than 10 times whithin 5 minutes."
		} else {
			userClickTracker[userID] = append(userClickTracker[userID], currentTime)

			// deduplicate click
			if !checkDuplicateClick(clickParams.ImpressionID) {
				clickTime := time.Now()

				// click expiration and daskhor
				eventsMutex.Lock()
				for i, event := range impressionEvents {
					if event.ID == clickParams.ImpressionID {
						// timeDiff := clickTime.Sub(event.Time).Seconds()
						timeDiff := clickParams.ExpTime.Sub(clickTime).Seconds()
						if timeDiff > 10 && timeDiff < 30 {
							impressionEvents[i].IsClicked = true
							impressionEvents[i].ClickID = clickParams.ID
							var updateApi = common.EventServiceApiModel{
								Time:         clickTime,
								UserID:       uuid.MustParse(userID),
								PubId:        clickParams.Pid,
								AdId:         clickParams.AdId,
								IsClicked:    true,
								ClickID:      clickParams.ID,
								ImpressionID: clickParams.ImpressionID,
							}
							ch <- updateApi
						} else {
							responseBody["click time error"] = "Click time is not in the valid range"
						}
						break
					}
				}
				eventsMutex.Unlock()
			} else {
				responseBody["duplicate error"] = "Duplicate click"
			}
		}
		var ad common.Ad
		result := Db.First(&ad, clickParams.AdId)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "adNum not found"})
			return
		}
		responseBody["AdURL"] = ad.Url
		c.JSON(http.StatusOK, responseBody)
	}
}

func checkDuplicateClick(impressionId uuid.UUID) bool {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	for _, event := range impressionEvents {
		if event.ID == impressionId && event.IsClicked {
			return true
		}
	}
	return false
}

func filterRecentClicks(clicks []time.Time, cutoff time.Time) []time.Time {
	recent := clicks[:0]
	for _, click := range clicks {
		if click.After(cutoff) {
			recent = append(recent, click)
		}
	}
	return recent
}

func panelApiCall(ch chan common.EventServiceApiModel) {
	for {
		select {
		case event := <-ch:
			fmt.Printf("channel size : %v", len(ch))

			jsonData, err := json.Marshal(event)
			if err != nil {
				fmt.Printf("can not umarshal event %s\n", err)
			}
			fmt.Printf("event %s\n", jsonData)
			http.Post("http://localhost:8085/eventservice", "application/json", bytes.NewReader(jsonData))
		default:
		}
	}
}

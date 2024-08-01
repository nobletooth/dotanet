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

		// decryption
		encryptedClickParams := c.Param("encryptedClickParams")
		decryptedClickParams, _ := decrypt(encryptedClickParams)
		fmt.Println(decryptedClickParams)
		var clickParams common.ClickedEvent
		json.Unmarshal(decryptedClickParams, &clickParams)
		adID := clickParams.AdId
		pubID := clickParams.Pid
		clickId := clickParams.ID
		impressionId := clickParams.ImpressionID

		// deduplicate click
		if !checkDuplicateClick(impressionId) {
			clickTime := time.Now()

			// click expiration and daskhor
			eventsMutex.Lock()
			for i, event := range impressionEvents {
				if event.ImpressionID == impressionId {
					timeDiff := clickTime.Sub(event.Time).Seconds()
					if timeDiff > 2 && timeDiff < 30 {
						impressionEvents[i].IsClicked = true
						impressionEvents[i].ClickID = clickId
						var updateApi = common.EventServiceApiModel{
							Time:         clickTime,
							PubId:        pubID,
							AdId:         adID,
							IsClicked:    true,
							ClickID:      clickId,
							ImpressionID: impressionId,
						}
						ch <- updateApi
					} else {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Click time is not in the valid range"})
					}
					break
				}
			}
			eventsMutex.Unlock()
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Duplicate click"})
			return
		}
		var ad common.Ad
		result := Db.First(&ad, adID)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "adNum not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"AdURL": ad.Url})
	}
}

func checkDuplicateClick(impressionId uuid.UUID) bool {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	for _, event := range impressionEvents {
		if event.ImpressionID == impressionId && event.IsClicked {
			return true
		}
	}
	return false
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

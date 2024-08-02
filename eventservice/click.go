package main

import (
	"bytes"
	"common"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func clickHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err    error
			advNum uint64
		)

		adv := c.Param("adv") // should decrypt adv and pub.
		if advNum, err = strconv.ParseUint(adv, 10, 32); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		advNum32 := uint(advNum)

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

		clickId := uuid.MustParse(c.Param("clickid"))
		impressionId := uuid.MustParse(c.Param("impressionid"))

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
							PubId:        pubInt,
							AdId:         adInt,
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
		result := Db.First(&ad, advNum32)
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

func panelApiCall(ch chan common.EventServiceApiModel, producer *kafka.Producer) {
	for {
		select {
		case event := <-ch:
			fmt.Printf("channel size : %v", len(ch))

			jsonData, err := json.Marshal(event)
			if err != nil {
				fmt.Printf("can not umarshal event %s\n", err)
			}
			panelUrl := *Panelserviceurl + "/eventservice"
			resp, err := http.Post(panelUrl, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Errorf("Error making POST request: %s\n", err)
			}
			//defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Received non-OK response status: %s\n", resp.Status)
			}

			err = producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &[]string{"clickview"}[0], Partition: kafka.PartitionAny},
				Value:          jsonData,
			}, nil)

			if err != nil {
				fmt.Printf("Error Posting to kafka %s\n", err)
			}
		default:
		}
	}
}

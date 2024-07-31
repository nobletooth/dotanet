package main

import (
	"common"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var processedClicks = make(map[string]bool)
var mu sync.Mutex

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
		time, err := time.Parse(time.RFC3339, c.Param("time"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time"})
			return
		}
		//in this part I want to check  don't double-click
		key := fmt.Sprintf("%d:%s", advNum32, pub)
		mu.Lock()
		if _, found := processedClicks[key]; !found {
			processedClicks[key] = true
			mu.Unlock()
			var updateApi = common.EventServiceApiModel{Time: time,
				PubId: pubInt, AdId: adInt, IsClicked: true, ClickID: clickId,
				ImpressionID: impressionId}
			ch <- updateApi
		} else {
			mu.Unlock()
		}

		//this is the end-of-it

		var ad common.Ad
		result := Db.First(&ad, advNum32)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "adNum not found"})
			return
		}
		//c.Redirect(http.StatusOK, "ad.Url")
		c.JSON(http.StatusOK, gin.H{"AdURL": ad.Url})
	}
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
			err = p.Produce(&kafka.Message{
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

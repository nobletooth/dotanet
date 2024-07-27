package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
)

func clickHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err    error
			advNum uint64
		)
		var clickTime = time.Now()
		adv := c.Param("adv") // should decrypt adv and pub.
		if advNum, err = strconv.ParseUint(adv, 10, 32); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		advNum32 := uint(advNum)
		pub := c.Param("pub") // should decrypt adv and pub.
		var updateApi = common.EventServiceApiModel{Time: clickTime, PubId: pub, AdId: adv, IsClicked: true}
		ch <- updateApi
		var ad common.Ad
		result := Db.First(&ad, advNum32)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "advNum not found"})
		}
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
				fmt.Errorf("error : " + err.Error())
			}
			resp, err := http.Post("http://localhost:8080/eventservice", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Errorf("Error making POST request: %s\n", err)
			}
			//defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Received non-OK response status: %s\n", resp.Status)
			}
		default:
		}
	}
	//jsonData, err := json.Marshal(<-ch)
	//if err != nil {
	//	fmt.Errorf("error : " + err.Error())
	//}
	//resp, err := http.Post(*EventservicePort+"/eventservice", "application/json", bytes.NewBuffer(jsonData))
	//_ = resp
}

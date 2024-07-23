package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type click struct {
	gorm.Model
	clickTime time.Time
	AdInfoId  uint `gorm:"foreignkey:AdId"`
}

type AdInfo struct {
	Id           uint    `json:"id"`
	Title        string  `json:"title"`
	Image        string  `json:"image"`
	Price        float64 `json:"price"`
	Status       bool    `json:"status"`
	Impressions  int     `json:"impressions"`
	Url          string  `json:"url"`
	AdvertiserId uint64  `json:"advertiserId"`
}

type updateApi struct {
	time      time.Time
	pubId     string
	adId      string
	isClicked bool
}

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
		var updateApi = updateApi{time: clickTime, pubId: pub, adId: adv, isClicked: true}
		ch <- updateApi
		var ad AdInfo
		result := Db.First(&ad, advNum32)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "advNum not found"})
		}

		c.Redirect(http.StatusMovedPermanently, ad.Url)

	}
}

func panelApiCall(ch chan updateApi) {
	var data = <-ch
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("error : " + err.Error())
	}
	resp, err := http.Post("http://localhost:8080/eventservice", "application/json", bytes.NewBuffer(jsonData))
	if resp.Status != "200 OK" {
		ch <- data
	}
}

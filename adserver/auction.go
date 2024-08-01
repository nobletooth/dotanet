package main

import (
	"common"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"slices"
	_ "slices"
	"strconv"
	"time"
)

func GetImage(adID uint) (string, error) {
	url := fmt.Sprintf("%v/ads/%d/picture", *PanelUrlPicUrl, adID)
	return url, nil
}

func GetAdHandler(c *gin.Context) {
	pubID := c.Param("pubID")

	if len(allAds) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No ads available"})
		return
	}
	pubidint, _ := strconv.Atoi(pubID)

	var newAds []common.AdWithMetrics
	var experiencedAds []common.AdWithMetrics
	var ctrPrices []float64
	var acceptbleadd []common.AdWithMetrics
	for _, firstad := range allAds {
		if slices.Contains(firstad.PreferdPubID, uint(pubidint)) || firstad.PreferdPubID[0] == 0 {
			acceptbleadd = append(acceptbleadd, firstad)
		}

	}

	for _, ad := range acceptbleadd {
		if ad.ImpressionCount < *NewAdImpressionThreshold {
			newAds = append(newAds, ad)
		} else {
			ctrPrice := float64(ad.ClickCount) / float64(ad.ImpressionCount) * ad.Price
			ctrPrices = append(ctrPrices, ctrPrice)
			experiencedAds = append(experiencedAds, ad)
		}
	}

	var selectedNewAd common.AdWithMetrics
	if len(newAds) > 0 {
		rand.Seed(time.Now().UnixNano())
		selectedNewAd = newAds[rand.Intn(len(newAds))]
	}

	totalScore := 0.0
	for _, ctrPrice := range ctrPrices {
		totalScore += ctrPrice
	}

	randomPoint := rand.Float64() * totalScore
	currentSum := 0.0
	var selectedExperiencedAd common.AdWithMetrics

	for i, ad := range experiencedAds {
		currentSum += ctrPrices[i]
		if randomPoint <= currentSum {
			selectedExperiencedAd = ad
			break
		}
	}

	var finalAd common.AdWithMetrics
	if (rand.Float64() < *NewAdSelectionProbability && selectedNewAd.Id != 0) || selectedExperiencedAd.Id == 0 {
		finalAd = selectedNewAd
	} else {
		finalAd = selectedExperiencedAd
	}

	sendAdResponse(c, finalAd, pubID)
}

func sendAdResponse(c *gin.Context, ad common.AdWithMetrics, pubID string) {
	fmt.Printf("adID: %d,ad title:%s,ad price:%f", ad.Id, ad.Title, ad.Price)
	imageDataurl, err := GetImage(ad.AdInfo.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image"})
		return
	}
	publisherID, err := strconv.ParseUint(pubID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publisher ID"})
		return
	}
	response := gin.H{
		"Title":          ad.Title,
		"ImageData":      imageDataurl,
		"ClicksURL":      fmt.Sprintf("%v/click/%d/%d", *EventServiceUrl, ad.Id, publisherID),
		"ImpressionsURL": fmt.Sprintf("%v/impression/%d/%d", *EventServiceUrl, ad.Id, publisherID),
	}

	c.JSON(http.StatusOK, response)
}

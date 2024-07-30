package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nobletooth/dotanet/common"
)

func GetImage(adID uint) (string, error) {
	url := fmt.Sprintf("%v/ads/%d/picture", *PanelUrl, adID)
	return url, nil
}

func GetAdHandler(c *gin.Context) {
	pubID := c.Param("pubID")

	if len(allAds) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No ads available"})
		return
	}

	var newAds []common.AdWithMetrics
	var experiencedAds []common.AdWithMetrics
	var ctrPrices []float64

	for _, ad := range allAds {
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

	var impression = common.ViewedEvent{
		ID:   uuid.New(),
		Time: time.Now(),
	}

	var click = common.ClickedEvent{
		ID:           uuid.New(),
		Time:         time.Now(),
		ImpressionID: impression.ID,
	}

	response := gin.H{
		"Title":          ad.Title,
		"ImageData":      imageDataurl,
		"ClicksURL":      fmt.Sprintf("%v/click/%d/%d/%v/%v/%v", *EventServiceUrl, ad.Id, publisherID, click.ID, click.ImpressionID, click.Time),
		"ImpressionsURL": fmt.Sprintf("%v/impression/%d/%d/%v/%v", *EventServiceUrl, ad.Id, publisherID, impression.ID, impression.Time),
	}

	c.JSON(http.StatusOK, response)
}

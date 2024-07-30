package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
)

func GetImage(adID uint) (string, error) {
	url := fmt.Sprintf("%v/ads/%d/picture", *PanelUrl, adID)
	return url, nil
}

func GetAdHandler(c *gin.Context) {
	pubID := c.Param("pubID")

	if len(allAds) == 0 {
		sendAdResponse(c, nil, pubID)
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

	rand.Seed(time.Now().UnixNano())
	selectNewAd := rand.Float64() < *NewAdSelectionProbability

	var finalAd *common.AdWithMetrics
	if selectNewAd && len(newAds) > 0 {
		rand.Seed(time.Now().UnixNano())
		selectedAd := newAds[rand.Intn(len(newAds))]
		finalAd = &selectedAd
	} else if len(experiencedAds) > 0 {
		totalScore := 0.0
		for _, ctrPrice := range ctrPrices {
			totalScore += ctrPrice
		}

		rand.Seed(time.Now().UnixNano())
		randomPoint := rand.Float64() * totalScore
		currentSum := 0.0
		for i, ad := range experiencedAds {
			currentSum += ctrPrices[i]
			if randomPoint <= currentSum {
				finalAd = &ad
				break
			}
		}
	} else {
		finalAd = nil
	}

	sendAdResponse(c, finalAd, pubID)
}

func sendAdResponse(c *gin.Context, ad *common.AdWithMetrics, pubID string) {
	if ad == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No ads available"})
		return
	}

	fmt.Printf("adID: %d, ad title: %s, ad price: %f\n", ad.Id, ad.Title, ad.Price)
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

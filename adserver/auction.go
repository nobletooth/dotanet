package main

import (
	"common"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"slices"

	_ "slices"
	"strconv"
	"time"
)


var RandomGenerator randomGenerator = &defaultRandomGenerator{}

type randomGenerator interface {
	Float64() float64
	Intn(n int) int
}

type defaultRandomGenerator struct{}

func (r *defaultRandomGenerator) Float64() float64 {
	return rand.Float64()
}

func (r *defaultRandomGenerator) Intn(n int) int {
	return rand.Intn(n)
}

func GetImage(adID uint) (string, error) {
	url := fmt.Sprintf("%v/ads/%d/picture", *PanelUrlPicUrl, adID)
	return url, nil
}

func GetAdHandler(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())

	pubID := c.Param("pubID")

	if len(allAds) == 0 {
		sendAdResponse(c, nil, pubID)
		return
	}
	pubidint, _ := strconv.Atoi(pubID)

	var newAds []common.AdWithMetrics
	var experiencedAds []common.AdWithMetrics
	var ctrPrices []float64
	var acceptbleadd []common.AdWithMetrics
	for _, firstad := range allAds {
		if len(firstad.PreferdPubID) == 0 {
			continue
		} else if slices.Contains(firstad.PreferdPubID, uint(pubidint)) || firstad.PreferdPubID[0] == 0 {
			acceptbleadd = append(acceptbleadd, firstad)
		}

	}
	fmt.Println(acceptbleadd)

	for _, ad := range acceptbleadd {
		if ad.ImpressionCount < *NewAdImpressionThreshold {
			newAds = append(newAds, ad)
		} else {
			if ad.Price > 0 {
				ctrPrice := float64(ad.ClickCount) / float64(ad.ImpressionCount) * ad.Price
				ctrPrices = append(ctrPrices, ctrPrice)
				experiencedAds = append(experiencedAds, ad)
			}
		}
	}

	selectNewAd := RandomGenerator.Float64() < *NewAdSelectionProbability

	var finalAd *common.AdWithMetrics
	if selectNewAd && len(newAds) > 0 {
		selectedAd := newAds[RandomGenerator.Intn(len(newAds))]
		finalAd = &selectedAd
	} else if len(experiencedAds) > 0 {
		totalScore := 0.0
		for _, ctrPrice := range ctrPrices {
			totalScore += ctrPrice
		}

		randomPoint := RandomGenerator.Float64() * totalScore
		currentSum := 0.0
		for i, ad := range experiencedAds {
			currentSum += ctrPrices[i]
			if randomPoint <= currentSum {
				finalAd = &ad
				break
			}
		}
	} else if len(newAds) > 0 {
		// Handle fall back
		selectedAd := newAds[RandomGenerator.Intn(len(newAds))]
		finalAd = &selectedAd
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
		fmt.Printf("\nFailed to get image\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image"})
		return
	}
	publisherID, err := strconv.ParseUint(pubID, 10, 64)
	if err != nil {
		fmt.Printf("\nFailed to get publisherID\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publisher ID"})
		return
	}

	var impression = common.ViewedEvent{
		ID: uuid.New(),
	}

	var click = common.ClickedEvent{
		ID:           uuid.New(),
		ImpressionID: impression.ID,
	}

	response := gin.H{
		"Title":          ad.Title,
		"ImageData":      imageDataurl,
		"ClicksURL":      fmt.Sprintf("%v/click/%d/%d/%v/%v", *EventServiceUrl, ad.Id, publisherID, click.ID, click.ImpressionID),
		"ImpressionsURL": fmt.Sprintf("%v/impression/%d/%d/%v", *EventServiceUrl, ad.Id, publisherID, impression.ID),
	}
	c.JSON(http.StatusOK, response)
}

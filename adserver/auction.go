package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	NewAdImpressionThreshold = flag.Int("newAdTreshold", 5, "Impression threshold for considering an ad as new")
	NewAdSelectionProbability = flag.Float64("newAdProb", 0.25, "Probability of selecting a new ad")
	ExperiencedAdSelectionProbability = flag.Float64("expAdProb", 0.75, "Probability of selecting a exprienced ad")
)

func GetImagePath(adID uint) (string, error) {
	url := fmt.Sprintf("http://localhost:8080/ads/%d/pictures", adID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get image path, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
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
		if ad.ImpressionCount < NewAdImpressionThreshold {
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
	if rand.Float64() < NewAdSelectionProbability && selectedNewAd.Id != 0 {
		finalAd = selectedNewAd
	} else {
		finalAd = selectedExperiencedAd
	}

	sendAdResponse(c, finalAd, pubID)
}

func sendAdResponse(c *gin.Context, ad common.AdWithMetrics, pubID string) {
	imagePath, err := GetImagePath(ad.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image path"})
		return
	}

	publisherID, err := strconv.ParseUint(pubID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publisher ID"})
		return
	}

	response := gin.H{
		"Title":          ad.Title,
		"ImagePath":      imagePath,
		"ClicksURL":      fmt.Sprintf("/click/%d/%d", ad.Id, publisherID),
		"ImpressionsURL": fmt.Sprintf("/impression/%d/%d", ad.Id, publisherID),
	}

	c.JSON(http.StatusOK, response)
}

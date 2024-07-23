package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
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

	sort.Slice(allAds, func(i, j int) bool {
		return allAds[i].Price*float64(allAds[i].Clicks)/float64(allAds[i].Impressions) >
			allAds[j].Price*float64(allAds[j].Clicks)/float64(allAds[j].Impressions)
	})

	bestAd := allAds[0]

	publisherID, err := strconv.ParseUint(pubID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publisher ID"})
		return
	}

	imagePath, err := GetImagePath(bestAd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image path"})
		return
	}

	response := gin.H{
		"Title":          bestAd.Title,
		"ImagePath":      imagePath,
		"ClicksURL":      fmt.Sprintf("/click/%d/%d", bestAd.AdvertiserId, publisherID),
		"ImpressionsURL": fmt.Sprintf("/impression/%d/%d", bestAd.AdvertiserId, publisherID),
	}

	c.JSON(http.StatusOK, response)
}

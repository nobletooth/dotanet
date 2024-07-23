package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	response := gin.H{
		"Title":          bestAd.Title,
		"ImagePath":      bestAd.Image,
		"ClicksURL":      fmt.Sprintf("/click/%d/%d", bestAd.AdvertiserId, publisherID),
		"ImpressionsURL": fmt.Sprintf("/impression/%d/%d", bestAd.AdvertiserId, publisherID),
	}

	c.JSON(http.StatusOK, response)
}

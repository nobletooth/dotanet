package main

func GetAdHandler(c *gin.Context) {
	if len(allAds) == 0{
		c.JSON(http.StatusNotFound, gin.H{"error": "No ads available"})
	}

	sort.Slice(allAds, func(i, j int) bool {
		return allAds[i].Price * float64(allAds[i].Clicks) / float64(allAds[i].Impressions) >
			allAds[j].Price * float64(allAds[j].Clicks) / float64(allAds[j].Impressions)
	})

	bestAd := allAds[0]

	c.JSON(http.StatusOK, bestAd)
}
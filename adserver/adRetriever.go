package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/nobletooth/dotanet/common"
)

func GetAdsListPeriodically() []common.AdWithMetrics {
	fetchads := func() {
		response, err := http.Get(*PanelUrlGetAllAds + "/ads/list/")
		if err != nil {
			log.Println("Error fetching ads list:", err)
			return
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			var result struct {
				Ads []common.AdWithMetrics `json:"ads"`
			}

			err = json.NewDecoder(response.Body).Decode(&result)
			if err != nil {
				log.Println("Error decoding response body:", err)
			} else {
				log.Println("Ads list fetched successfully")
				allAds = ReturnAllAds(result.Ads)
			}
		} else {
			log.Println("Failed to fetch ads list, status code:", response.StatusCode)
		}
	}

	fetchads()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fetchads()

		}
	}
}

func ReturnAllAds(ads []common.AdWithMetrics) []common.AdWithMetrics {
	if ads == nil {
		log.Println("No ads found")
		return []common.AdWithMetrics{}
	}
	for _, ad := range ads {
		log.Printf("Ad ID: %d, Title: %s", ad.AdInfo.Id, ad.AdInfo.Title)
	}
	return ads
}

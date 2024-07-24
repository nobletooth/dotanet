package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/nobletooth/dotanet/common"
)

func GetAdsListPeriodically() []common.AdInfo {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			response, err := http.Get("http://localhost:" + PanelPort + "/ads/list/")
			if err != nil {
				log.Println("Error fetching ads list:", err)
				continue
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				var ads []common.AdInfo
				err = json.NewDecoder(response.Body).Decode(&ads)
				if err != nil {
					log.Println("Error decoding response body:", err)
				} else {
					log.Println("Ads list fetched successfully")
					allAds = ReturnAllAds(ads)
				}
			} else {
				log.Println("Failed to fetch ads list, status code:", response.StatusCode)
			}
		}
	}
}

func ReturnAllAds(ads []common.AdInfo) []common.AdInfo {
	if ads == nil {
		log.Println("No ads found")
		return []common.AdInfo{}
	}
	for _, ad := range ads {
		log.Printf("Ad ID: %d, Title: %s", ad.Id, ad.Title)
	}
	return ads
}

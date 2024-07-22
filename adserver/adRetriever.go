package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"example.com/dotanet/common"
)

func GetAdsListPeriodically(ads []common.Ad) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			response, err := http.Get("http://localhost:8080/ads/")
			if err != nil {
				log.Println("Error fetching ads list:", err)
				continue
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				err = json.NewDecoder(response.Body).Decode(&ads)
				if err != nil {
					log.Println("Error decoding response body:", err)
				} else {
					log.Println("Ads list fetched successfully")
					for _, ad := range ads {
						log.Printf("Ad ID: %d, Title: %s", ad.Id, ad.Title)
					}
				}
			} else {
				log.Println("Failed to fetch ads list, status code:", response.StatusCode)
			}
		}
	}
}

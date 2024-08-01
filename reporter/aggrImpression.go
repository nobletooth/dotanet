package main

import (
	"gorm.io/gorm"
	"time"
)

type aggrImpression struct {
	gorm.Model
	AdId            int
	ImpressionCount int
}

func addAggrImpressionDb() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		for key, value := range batchMapClick {
			if value > 5 {
				mu.Lock()
				aggrImpression := aggrImpression{AdId: key, ImpressionCount: value}
				DB.Create(&aggrImpression)
				batchMapClick[key] = 0
				mu.Unlock()
			}
		}
	}
}

package main

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type aggrImpression struct {
	gorm.Model
	AdId            int
	ImpressionCount int
}

func addAggrImpressionDb() {
	fmt.Printf("\nbefore ticker in addAggrImpressionDb\n")
	ticker := time.NewTicker(30 * time.Second)
	fmt.Printf("\nafter ticker in addAggrImpressionDb\n")
	for range ticker.C {
		fmt.Printf("\ninside for range before map checking in addAggrImpressionDb\n")
		for key, value := range batchMapImpression {
			fmt.Printf("\ninside for range after map checking in addAggrImpressionDb\n")
			if value > 0 {
				mu.Lock()
				aggrImpression := aggrImpression{AdId: key, ImpressionCount: value}
				DB.Create(&aggrImpression)
				fmt.Printf("\nsend impression data to db\n")
				batchMapImpression[key] = 0
				mu.Unlock()
			}
		}
	}
}

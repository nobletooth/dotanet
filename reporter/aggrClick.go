package main

import (
	"gorm.io/gorm"
	"time"
)

type aggrClick struct {
	gorm.Model
	AdId       int
	ClickCount int
}

func addAggrClickDb() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		for key, value := range batchMapClick {
			if value > 5 {
				mu.Lock()
				aggrClick := aggrClick{AdId: key, ClickCount: value}
				DB.Create(&aggrClick)
				batchMapClick[key] = 0
				mu.Unlock()
			}
		}
	}
}

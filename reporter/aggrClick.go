package main

import (
	"fmt"
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
			if value > 0 {
				mu.Lock()
				aggrClick := aggrClick{AdId: key, ClickCount: value}
				DB.Create(&aggrClick)
				fmt.Printf("\nsend click data to db\n")
				batchMapClick[key] = 0
				mu.Unlock()
			}
		}
	}
}

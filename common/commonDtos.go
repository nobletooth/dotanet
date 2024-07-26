package common

import (
	"time"
)

type AdInfo struct {
	Id           uint    `json:"id"`
	Title        string  `json:"title"`
	Image        string  `json:"image"`
	Price        float64 `json:"price"`
	Status       bool    `json:"status"`
	Impressions  int     `json:"impressions"`
	Url          string  `json:"url"`
	AdvertiserId uint64  `json:"advertiserId"`
}

type ClickedEvent struct {
	ID   uint      `gorm:"primaryKey"`
	Pid  string    `gorm:"index"`
	AdId string    `gorm:"index"`
	Time time.Time `gorm:"index"`
}

type ViewedEvent struct {
	ID   uint      `gorm:"primaryKey"`
	Pid  string    `gorm:"index"`
	AdId string    `gorm:"index"`
	Time time.Time `gorm:"index"`
}

type EventServiceApiModel struct {
	Time      time.Time
	PubId     string
	AdId      string
	IsClicked bool
}

type AdWithMetrics struct {
	AdInfo
	ClickCount      int64 `json:"clickCount"`
	ImpressionCount int64 `json:"impressionCount"`
}

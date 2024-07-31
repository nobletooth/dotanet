package common

import (
	"time"

	"github.com/google/uuid"
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
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Pid          int       `gorm:"index"`
	AdId         int       `gorm:"index"`
	Time         time.Time `gorm:"index"`
	ImpressionID uuid.UUID `gorm:"type:uuid;unique;foreignKey:ID"`
}

type ViewedEvent struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Pid  int       `gorm:"index"`
	AdId int       `gorm:"index"`
	Time time.Time `gorm:"index"`
}

type EventServiceApiModel struct {
	Time         time.Time
	PubId        int
	AdId         int
	IsClicked    bool
	ClickID      uuid.UUID
	ImpressionID uuid.UUID
}

type AdWithMetrics struct {
	AdInfo
	ClickCount      int64 `json:"clickCount"`
	ImpressionCount int64 `json:"impressionCount"`
}

type Ad struct {
	Id           uint    `gorm:"column:id;primary_key"`
	Title        string  `gorm:"column:title"`
	Image        string  `gorm:"column:image"`
	Price        float64 `gorm:"column:price"`
	Status       bool    `gorm:"column:status"`
	Clicks       int     `gorm:"column:clicks"`
	Impressions  int     `gorm:"column:impressions"`
	Url          string  `gorm:"column:url"`
	AdvertiserId uint64  `gorm:"foreignKey:AdvertiserId"`
}

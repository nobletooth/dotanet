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

type UrlClickParameters struct {
	ID           uuid.UUID `json:"id"`
	Pid          int       `json:"pid"`
	AdId         int       `json:"adid"`
	Time         time.Time `json:"time"`
	ImpressionID uuid.UUID `json:"impressionid"`
	ExpTime      time.Time `json:"exptime"`
}

type UrlImpressionParameters struct {
	ID         uuid.UUID `json:"id"`
	Pid        int       `json:"pid"`
	AdId       int       `json:"adid"`
	Time       time.Time `json:"time"`
	IsClicked  bool      `json:"isclicked"`
	LoadAdTime time.Time `json:"loadadtime"`
	ClickID    uuid.UUID `json:"clickid"`
}

type ClickedEvent struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID `gorm:"type:uuid"`
	Pid          int       `gorm:"index"`
	AdId         int       `gorm:"index"`
	Time         time.Time `gorm:"index"`
	ImpressionID uuid.UUID `gorm:"type:uuid;unique;foreignKey:ID"`
}

type ViewedEvent struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid"`
	Pid    int       `gorm:"index"`
	AdId   int       `gorm:"index"`
	Time   time.Time `gorm:"index"`
}

type EventServiceApiModel struct {
	UserID       uuid.UUID
	Time         time.Time
	PubId        int
	AdId         int
	IsClicked    bool
	ClickID      uuid.UUID
	ImpressionID uuid.UUID
}

type AdWithMetrics struct {
	AdInfo
	ClickCount      int64  `json:"clickCount"`
	ImpressionCount int64  `json:"impressionCount"`
	PreferdPubID    []uint `json:"preferd-pub-id"`
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

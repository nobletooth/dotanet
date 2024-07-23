package common

import "time"

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

type click struct {
	Id       uint
	time     time.Time
	AdInfoId uint
}

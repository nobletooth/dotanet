package main

type AdInfo struct {
	Id           uint    `json:"id"`
	Title        string  `json:"title"`
	Image        string  `json:"image"`
	Price        float64 `json:"price"`
	Status       bool    `json:"status"`
	Clicks       int     `json:"clicks"`
	Impressions  int     `json:"impressions"`
	Url          string  `json:"url"`
	AdvertiserId uint64  `json:"advertiserId"`
}

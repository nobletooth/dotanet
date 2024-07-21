package advertiser

type Ad struct {
	Id          uint    `json:"id"`
	Title       string  `json:"title`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Status      bool    `json:"status"`
	Url         string  `json:"url"`
	Clicks      int     `json:"clicks"`
	Impressions int     `json:"impressions"`
}

func CTR(ad *Ad) float64 {
	if ad.Impressions == 0 {
		return 0
	}
	return float64(ad.Clicks) / float64(ad.Impressions)
}
func CostCalculator(ad *Ad) float64 {
	return float64(ad.Clicks) * ad.Price
}

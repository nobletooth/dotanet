package advertiser

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdInfo struct {
	Id           uint    `json:"id"`
	Title        string  `json:"title"`
	Image        string  `json:"image"`
	Price        float64 `json:"price"`
	Status       bool    `json:"status"`
	Url          string  `json:"url"`
	AdvertiserId uint    `gorm:"foreignKey:advertiserid"`
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
	AdvertiserId uint    `gorm:"foreignKey:advertiserid"`
}

type AdService struct {
	db *Database
}

func NewAdService(db *Database) AdService {
	return AdService{db: db}
}

func (ad *AdInfo) CreateAd(db *Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		if ad.Title == "" || ad.Image == "" || ad.Price == 0 || ad.Url == "" || ad.AdvertiserId == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}
		ad.Status = true
		entity := MapEntity(*ad)

		if err := db.DB.Create(&entity).Error; err != nil {
			log.Printf("Creating ad failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Creating ad failed"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Ad created successfully", "ad": entity})
	}

}

func MapEntity(ad AdInfo) Ad {
	return Ad{
		Id:           ad.Id,
		Title:        ad.Title,
		Image:        ad.Image,
		Price:        ad.Price,
		Status:       ad.Status,
		Url:          ad.Url,
		AdvertiserId: ad.AdvertiserId,
		Clicks:       0,
		Impressions:  0,
	}
}

func CTR(entity Ad) float64 {
	if entity.Impressions == 0 {
		return 0
	}
	return float64(entity.Clicks) / float64(entity.Impressions)
}
func CostCalculator(entity Ad) float64 {
	return float64(entity.Clicks) * entity.Price

}

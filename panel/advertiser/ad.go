package advertiser

import (
    "dotanet/db"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

type Ad struct {
    Id     uint    `json:"id"`
    Title  string  `json:"title"`
    Image  string  `json:"image"`
    Price  float64 `json:"price"`
    Status bool    `json:"status"`
    Url    string  `json:"url"`
}

type AdEntity struct {
    Id          uint    `gorm:"column:id;primary_key"`
    Title       string  `gorm:"column:title"`
    Image       string  `gorm:"column:image"`
    Price       float64 `gorm:"column:price"`
    Status      bool    `gorm:"column:status"`
    Clicks      int     `gorm:"column:clicks"`
    Impressions int     `gorm:"column:impressions"`
    Url         string  `gorm:"column:url"`
}

func (ad *Ad) CreateAd(db *db.Database) gin.HandlerFunc {

    return func(c *gin.Context) {
        if ad.Title == "" || ad.Image == "" || ad.Price == 0 || ad.Url == "" {
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

func MapEntity(ad Ad) AdEntity {
    return AdEntity{
        Id:          ad.Id,
        Title:       ad.Title,
        Image:       ad.Image,
        Price:       ad.Price,
        Status:      ad.Status,
        Url:         ad.Url,
        Clicks:      0,
        Impressions: 0,
    }
}

func CTR(entity AdEntity) float64 {
    if entity.Impressions == 0 {
        return 0
    }
    return float64(entity.Clicks) / float64(entity.Impressions)
}
func CostCalculator(entity AdEntity) float64 {
    return float64(entity.Clicks) * entity.Price
}

func (ad *Ad) DeleteAd() {

}


package advertiser

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

// type Ad struct {
// 	Id           uint    `json:"id"`
// 	Title        string  `json:"title"`
// 	Image        string  `json:"image"`
// 	Price        float64 `json:"price"`
// 	Status       bool   `json:"status"`
// 	Clicks       int     `json:"clicks"`
// 	Impressions  int     `json:"impressions"`
// 	Url          string  `json:"url"`
// 	AdvertiserId uint    `gorm:"foreignKey:advertiserid"`
// }

func CreateAd(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form Ad
		if err := c.Bind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if form.Title == "" || form.Image == "" || form.Price == 0 || form.Url == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		form.Status = true
		if err := db.Create(&form).Error; err != nil {
			log.Printf("Creating ad failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Creating ad failed"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Ad created successfully", "ad": form})
	}
}

func DeleteAd(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adID := c.Param("id")
		var ad Ad
		if err := db.First(&ad, adID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Ad not found"})
			} else {
				log.Printf("Error finding ad: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding ad"})
			}
			return
		}
		if err := db.Delete(&ad).Error; err != nil {
			log.Printf("Deleting ad failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Deleting ad failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Ad deleted successfully"})
	}
}
func UpdateAd(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adID := c.Param("id")
		var ad Ad
		if err := db.First(&ad, adID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Ad not found"})
			} else {
				log.Printf("Error finding ad: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding ad"})
			}
			return
		}

		var updatedForm Ad
		if err := c.Bind(&updatedForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Model(&ad).Updates(updatedForm).Error; err != nil {
			log.Printf("Updating ad failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Updating ad failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Ad updated successfully", "ad": updatedForm})
	}
}

func RenderCreateAdForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create_ad.html", gin.H{})
}

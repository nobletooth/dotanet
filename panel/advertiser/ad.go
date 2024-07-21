package advertiser

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	AdvertiserId uint64  `gorm:"foreignKey:AdvertiserId"`
}

func CreateAdHandler(c *gin.Context) {
	var ad Ad
	ad.Title = c.PostForm("title")
	ad.Image = c.PostForm("image")
	ad.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	ad.Url = c.PostForm("url")
	ad.AdvertiserId, _ = strconv.ParseUint(c.PostForm("advertiser_id"), 10, 32)
	ad.Status = true

	if err := DB.Create(&ad).Error; err != nil {
		log.Printf("Creating ad failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Creating ad failed"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/advertisers/"+strconv.Itoa(int(ad.AdvertiserId))+"/ads")
}

func ListAdsByAdvertiserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index", gin.H{"error": "Invalid advertiser ID"})
		return
	}

	ads, err := ListAdsByAdvertiser(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "advertiser_ads", gin.H{"Ads": ads})
}
func CreateAdForm(c *gin.Context) {
	advertiserID := c.Query("advertiser_id")
	c.HTML(http.StatusOK, "create_ad", gin.H{"AdvertiserID": advertiserID})
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

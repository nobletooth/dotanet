package advertiser

import (
	"errors"
	"github.com/nobletooth/dotanet/common"
	"github.com/nobletooth/dotanet/panel/database"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"

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
	ad.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	ad.Url = c.PostForm("url")
	ad.AdvertiserId, _ = strconv.ParseUint(c.PostForm("advertiser_id"), 10, 32)
	ad.Status = true

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	imagePath := "./image/" + file.Filename + ad.Url
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	ad.Image = imagePath

	if err := database.DB.Create(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Creating ad failed"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/advertisers/"+strconv.Itoa(int(ad.AdvertiserId))+"/ads")
}

func CreateAdForm(c *gin.Context) {
	advertiserID := c.Query("advertiser_id")
	id, err := strconv.Atoi(advertiserID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index", gin.H{"error": "Invalid advertiser ID"})
		return
	}

	advertiser, err := FindAdvertiserByID(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "create_ad", gin.H{"Advertiser": advertiser})
}

func FindAdvertiserByID(id uint) (Entity, error) {
	var entity Entity
	result := database.DB.First(&entity, id)
	if result.Error != nil {
		return Entity{}, result.Error
	}
	return entity, nil
}

func UpdateAdHandler(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "advertiser_ads", gin.H{"error": "Invalid ad ID"})
		return
	}

	var ad Ad
	if err := database.DB.First(&ad, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "advertiser_ads", gin.H{"error": "Ad not found"})
		return
	}

	ad.Title = c.PostForm("title")
	ad.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	ad.Url = c.PostForm("url")
	ad.Status = c.PostForm("status") == "on"

	if err := database.DB.Save(&ad).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "advertiser_ads", gin.H{"error": "Failed to update ad"})
		return
	}

	c.Redirect(http.StatusFound, "/advertisers/"+strconv.Itoa(int(ad.AdvertiserId))+"/ads")
}

func FindAdById(id int) (Ad, error) {
	var ad Ad
	err := database.DB.First(&ad, id).Error
	return ad, err

}

func LoadAdPictureHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "advertiser_ads", gin.H{"error": "Invalid ad ID"})
		return
	}
	id64 := uint(id)

	ad, err := FindAdById(int(id64))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HTML(http.StatusNotFound, "advertiser_ads", gin.H{"error": "Ad not found"})
		} else {
			c.HTML(http.StatusInternalServerError, "advertiser_ads", gin.H{"error": "Failed to retrieve ad"})
		}
		return
	}

	imageFilePath := ad.Image

	file, err := os.Open(imageFilePath)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "advertiser_ads", gin.H{"error": "Failed to open image file"})
		return
	}
	defer file.Close()
	filebytes, _ := os.ReadFile(imageFilePath)
	contentType := http.DetectContentType(filebytes)
	c.Header("Content-Type", contentType)
	c.File(imageFilePath)
}

func ListAllAds(c *gin.Context) {
	var ads []Ad
	result := database.DB.Find(&ads)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Loading ads failed"})
		return
	}
	startTime := time.Now().Add(time.Duration(-1) * time.Hour)
	endTime := time.Now()

	var adMetrics []common.AdWithMetrics

	for _, ad := range ads {
		var clickCount int64
		var impressionCount int64

		database.DB.Model(&common.ClickedEvent{}).
			Where("ad_id = ? AND time BETWEEN ? AND ?", ad.Id, startTime, endTime).
			Count(&clickCount)

		database.DB.Model(&common.ViewedEvent{}).
			Where("ad_id = ? AND time BETWEEN ? AND ?", ad.Id, startTime, endTime).
			Count(&impressionCount)

		adinfo := common.AdInfo{
			AdvertiserId: ad.AdvertiserId,
			Id:           ad.Id,
			Price:        ad.Price,
			Url:          ad.Url,
			Status:       ad.Status,
			Title:        ad.Title,
		}

		adMetrics = append(adMetrics, common.AdWithMetrics{
			AdInfo:          adinfo,
			ClickCount:      clickCount,
			ImpressionCount: impressionCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{"ads": adMetrics})
}

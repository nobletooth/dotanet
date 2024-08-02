package advertiser

import (
	"common"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/nobletooth/dotanet/panel/database"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Ad struct {
	Id           uint           `gorm:"column:id;primary_key"`
	Title        string         `gorm:"column:title"`
	Image        string         `gorm:"column:image"`
	Price        float64        `gorm:"column:price"`
	Status       bool           `gorm:"column:status"`
	Clicks       int            `gorm:"column:clicks"`
	Impressions  int            `gorm:"column:impressions"`
	Url          string         `gorm:"column:url"`
	Keyword      pq.StringArray `gorm:"type:text[]"`
	AdvertiserId uint64         `gorm:"foreignKey:AdvertiserId"`
	AdLimit      *float64       `gorm:"column:ad_limit"`
}

func handlekeyword(keyword string, clearkeyword bool, ad *Ad) {
	if clearkeyword == true {
		ad.Keyword = nil
	} else if keyword != "" {
		ad.Keyword = strings.Split(keyword, ",")
	}

}

func CreateAdHandler(c *gin.Context) {
	var ad Ad
	ad.Title = c.PostForm("title")
	ad.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	ad.Url = c.PostForm("url")
	keyword := c.PostForm("keyword")
	if keyword != "" {
		ad.Keyword = strings.Split(keyword, ",")
	}
	ad.AdvertiserId, _ = strconv.ParseUint(c.PostForm("advertiser_id"), 10, 32)
	ad.Status = true

	if adLimitStr := c.PostForm("ad_limit"); adLimitStr != "" {
		adLimit, _ := strconv.ParseFloat(adLimitStr, 64)
		ad.AdLimit = &adLimit
	} else {
		ad.AdLimit = nil
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	newFilename := fmt.Sprintf("%v-%s", ad.AdvertiserId, file.Filename)
	imageDir := filepath.Join(os.Getenv("HOME"), "Desktop", "dotanet", "panel", "image")
	imagePath := filepath.Join(imageDir, newFilename)

	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	ad.Image = newFilename
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

	advertiser, err := FindAdvertiserByID(uint64(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "create_ad", gin.H{"Advertiser": advertiser})
}

func FindAdvertiserByID(id uint64) (Entity, error) {
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
	clearkeyword := c.PostForm("clear-keyword-button") == "on"
	keyword := c.PostForm("keyword")
	handlekeyword(keyword, clearkeyword, &ad)

	if adLimitStr := c.PostForm("ad_limit"); adLimitStr != "" {
		adLimit, _ := strconv.ParseFloat(adLimitStr, 64)
		ad.AdLimit = &adLimit
	} else {
		ad.AdLimit = nil
	}

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

	fmt.Println("\n image file path: " + imageFilePath)
	file, err := os.Open(imageFilePath) //TODO
	if err != nil {
		fmt.Println("500")
		c.HTML(http.StatusInternalServerError, "advertiser_ads", gin.H{"error": "Failed to open image file"})
		return
	}
	defer file.Close()
	filebytes, _ := os.ReadFile(imageFilePath)
	contentType := http.DetectContentType(filebytes)
	c.Header("Content-Type", contentType)
	c.File(imageFilePath)
}

func GetPublisherIDsWithAllKeywords(keywords []string) ([]uint, error) {
	var publisherIDs []uint

	err := database.DB.Table("publishers").
		Select("publishers.id").
		Distinct().
		Joins("JOIN publisher_tokens ON publishers.id = publisher_tokens.publisher_id").
		Joins("JOIN publisher_token_associations ON publisher_tokens.id = publisher_token_associations.publisher_token_id").
		Joins("JOIN tokens ON publisher_token_associations.token_value = tokens.value").
		Where("tokens.value IN ?", keywords).
		Group("publishers.id").
		Having("COUNT(DISTINCT tokens.value) = ?", len(keywords)).
		Pluck("publishers.id", &publisherIDs).
		Error

	if err != nil {
		return nil, err
	}

	return publisherIDs, nil
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

		advertiser, err := FindAdvertiserByID(ad.AdvertiserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch advertiser data"})
			return
		}

		// Check if the ad status is active
		if !ad.Status {
			continue
		}

		// Advertiser credit check
		if float64(advertiser.Credit) <= ad.Price {
			continue
		}

		var totalClicks int64

		database.DB.Table("clicked_events").
			Where("ad_id = ?", ad.Id).
			Count(&totalClicks)

		//Ad limit check
		if ad.AdLimit != nil && totalClicks >= int64(*ad.AdLimit/ad.Price) {
			continue
		}

		var clickCount int64
		var impressionCount int64

		database.DB.Table("clicked_events").
			Joins("INNER JOIN viewed_events ON clicked_events.impression_id = viewed_events.id").
			Where("viewed_events.ad_id = ? AND viewed_events.time BETWEEN ? AND ?", ad.Id, startTime, endTime).
			Count(&clickCount)

		database.DB.Table("viewed_events").
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
		var pubids []uint
		if ad.Keyword == nil {
			pubids = []uint{0}
		} else {
			var err error
			pubids, err = GetPublisherIDsWithAllKeywords(ad.Keyword)
			if err != nil {
				log.Println(err)
			}
		}
		adMetrics = append(adMetrics, common.AdWithMetrics{
			AdInfo:          adinfo,
			ClickCount:      clickCount,
			ImpressionCount: impressionCount,
			PreferdPubID:    pubids,
		})
	}
	fmt.Printf("%v", adMetrics[0].PreferdPubID)

	c.JSON(http.StatusOK, gin.H{"ads": adMetrics})
}

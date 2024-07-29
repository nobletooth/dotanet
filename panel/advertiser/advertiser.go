package advertiser

import (
	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
	"github.com/nobletooth/dotanet/panel/database"
	"gorm.io/gorm"
	"math"
	_ "math"
	"net/http"
	"strconv"
	"time"
)

type Entity struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Name   string `gorm:"unique;not null"`
	Credit int    `gorm:"column:credit"`
	Ads    []Ad   `gorm:"foreignKey:AdvertiserId"`
}

type AdReport struct {
	Date        time.Time `json:"date"`
	Clicks      int64     `json:"clicks"`
	Impressions int64     `json:"impressions"`
	CTR         float64   `json:"ctr"`
	Spent       float64   `json:"spent"`
}

type Service interface {
	GetCreditOfAdvertiser(adId int) (Entity, error)
	CreateAdvertiserEntity(name string, credit int)
	ListAllAdvertiserEntities() []Entity
	FindAdvertiserByName(name string) (Entity, error)
	ListAdsByAdvertiser(advertiserId uint) ([]Ad, error)
}

func HandleAdvertiserCredit(tx *gorm.DB, ad Ad) error {
	creditDeduction := int(math.Ceil(ad.Price * 0.8))
	result := tx.Model(&Entity{}).Where("ID = ?", ad.AdvertiserId).Update("Credit", gorm.Expr("Credit - ?", creditDeduction))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetCreditOfAdvertiser(adId int) (Entity, error) {
	var entity Entity

	result := database.DB.First(&entity, adId)
	if result.Error != nil {
		return entity, result.Error
	}
	return entity, nil
}

func CreateAdvertiserEntity(name string, credit int) {
	entity := Entity{
		Name:   name,
		Credit: credit,
	}
	database.DB.Create(&entity)
}

func ListAllAdvertiserEntities() []Entity {
	var advertisers []Entity
	database.DB.Find(&advertisers)
	return advertisers
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

	advertiser, err := FindAdvertiserByID(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "advertiser_ads", gin.H{"Ads": ads, "Advertiser": advertiser})
}
func FindAdvertiserByName(name string) (Entity, error) {
	var entity Entity
	result := database.DB.Where("name = ?", name).First(&entity)
	if result.Error != nil {
		return Entity{}, result.Error
	}
	return entity, nil
}

func ListAdsByAdvertiser(advertiserId uint) ([]Ad, error) {
	var ads []Ad
	result := database.DB.Where("advertiser_id = ?", advertiserId).Find(&ads)
	return ads, result.Error
}

func ListAdvertisers(c *gin.Context) {
	advertisers := ListAllAdvertiserEntities()
	c.HTML(http.StatusOK, "advertisers", gin.H{"Advertisers": advertisers})
}
func NewAdvertiserForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create_advertiser", nil)
}

func CreateAdvertiser(c *gin.Context) {
	name := c.PostForm("name")
	credit, err := strconv.Atoi(c.PostForm("credit"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "create_advertiser", gin.H{"error": "Invalid credit amount"})
		return
	}
	CreateAdvertiserEntity(name, credit)
	c.Redirect(http.StatusSeeOther, "/")
}

func GetAdvertiserCredit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index", gin.H{"error": "Invalid advertiser ID"})
		return
	}

	credit, err := GetCreditOfAdvertiser(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "advertiser_credit", credit)
}

func EditAdForm(c *gin.Context) {
	adIDStr := c.Param("id")
	adID, err := strconv.Atoi(adIDStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index", gin.H{"error": "Invalid ad ID"})
		return
	}

	var ad Ad
	result := database.DB.First(&ad, adID)
	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "index", gin.H{"error": result.Error.Error()})
		return
	}

	c.HTML(http.StatusOK, "edit_ad", gin.H{"Ad": ad})
}

func GetAdvertiserAdReports(c *gin.Context) {
	adIDStr := c.Param("id")
	adID, err := strconv.Atoi(adIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ad ID"})
		return
	}

	now := time.Now()
	endDate := now.Truncate(time.Minute)
	startDate := endDate.Add(-1 * time.Hour)

	var reports []AdReport
	for date := startDate; !date.After(endDate); date = date.Add(time.Minute) {
		var clickCount, impressionCount int64
		var spent float64

		database.DB.Table("viewed_events").
			Select("COUNT(viewed_events.id) as impressionCount, COUNT(DISTINCT clicked_events.impression_id) as clickCount").
			Joins("LEFT JOIN clicked_events ON viewed_events.id = clicked_events.impression_id").
			Where("viewed_events.ad_id = ? AND viewed_events.time BETWEEN ? AND ?", ad.Id, startTime, endTime).
			Group("viewed_events.ad_id").
			Scan(&impressionCount, &clickCount)

		database.DB.Table("clicked_events").
			Select("SUM(ads.price)").
			Where("ad_id = ? AND time BETWEEN ? AND ?", adID, date, date.Add(time.Minute)).
			Scan(&spent)

		ctr := 0.0
		if impressionCount > 0 {
			ctr = float64(clickCount) / float64(impressionCount)
		}

		reports = append(reports, AdReport{
			Date:        date,
			Clicks:      clickCount,
			Impressions: impressionCount,
			CTR:         ctr,
			Spent:       spent,
		})
	}

	c.HTML(http.StatusOK, "ad_reports", gin.H{
		"Reports":   reports,
		"StartDate": startDate.Format("2006-01-02 15:04:05"),
		"EndDate":   endDate.Format("2006-01-02 15:04:05"),
		"AdID":      adID,
	})
}

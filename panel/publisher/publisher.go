package publisher

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/common"
	"github.com/nobletooth/dotanet/panel/advertiser"
	"github.com/nobletooth/dotanet/panel/database"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Publisher struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	URL         string `gorm:"column:url;unique;not null"` // Removed space after "url"
	Credit      int    `gorm:"column:credit"`
	Clicks      int    `gorm:"column:clicks"`
	Impressions int    `gorm:"column:impressions"`
}

type Report struct {
	Date        time.Time `json:"date"`
	Income      float64   `json:"income"`
	Clicks      int64     `json:"clicks"`
	Impressions int64     `json:"impressions"`
}

func HandlePublisherCredit(ad advertiser.Ad, pid int) error {
	creditAddition := int(ad.Price * 0.2)
	result := database.DB.Model(&Publisher{}).Where("ID = ?", pid).Update("Credit", gorm.Expr("Credit + ?", creditAddition))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func ListPublishers(c *gin.Context) {
	var publishers []Publisher
	if err := database.DB.Find(&publishers).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "publishers", gin.H{"error": "Failed to load publishers"})
		return
	}
	c.HTML(http.StatusOK, "publishers", gin.H{
		"Publishers": publishers,
	})
}

func NewPublisherForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create_publisher", nil)
}

func CreatePublisherHandler(c *gin.Context) {
	url := c.PostForm("url")
	publisher := Publisher{URL: url, Credit: 0}
	if err := database.DB.Create(&publisher).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "index", gin.H{"error": "Failed to create publisher"})
		return
	}

	c.Redirect(http.StatusFound, "/publishers")
}

func ViewPublisherHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index", gin.H{"error": "Invalid publisher ID"})
		return
	}

	var publisher Publisher
	if err := database.DB.First(&publisher, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "index", gin.H{"error": "Publisher not found"})
		return
	}

	script := "<script src=\"" + publisher.URL + ".com/script.js\"></script>"

	c.HTML(http.StatusOK, "view_publisher", gin.H{
		"Publisher": publisher,
		"Script":    script,
	})
}

func GetPublisherScript(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publisher ID"})
		return
	}

	var publisher Publisher
	if err := database.DB.First(&publisher, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Publisher not found"})
		return
	}

	scriptFilePath := "./publisher/script/" + "digikala" + ".js"
	scriptContent, err := os.ReadFile(scriptFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read script file"})
		return
	}

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	c.DataFromReader(200, int64(len(scriptContent)), "application/javascript", bytes.NewReader(scriptContent), map[string]string{})
}

func GetPublisherReports(c *gin.Context) {
	publisherIDStr := c.Param("id")
	publisherID, err := strconv.Atoi(publisherIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publisher ID"})
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}
	var reports []Report
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		var clickCount, impressionCount int64
		var income float64

		database.DB.Model(&common.ClickedEvent{}).
			Where("pid = ? AND time BETWEEN ? AND ?", publisherID, date, date.AddDate(0, 0, 1)).
			Count(&clickCount)

		database.DB.Model(&common.ViewedEvent{}).
			Where("pid = ? AND time BETWEEN ? AND ?", publisherID, date, date.AddDate(0, 0, 1)).
			Count(&impressionCount)

		database.DB.Model(&advertiser.Ad{}).
			Select("SUM(price * 0.2)").
			Joins("JOIN publishers ON publishers.id = ads.publisher_id").
			Where("publisher_id = ? AND ads.created_at BETWEEN ? AND ?", publisherID, date, date.AddDate(0, 0, 1)).
			Scan(&income)

		database.DB.Table("clicked_events").
			Select("SUM(price * 0.2)").
			Joins("JOIN ads ON clicked_events.ad_id = ads.id").
			Where("clicked_events.pid = ? AND clicked_events.time BETWEEN ? AND ?", publisherID, date, date.AddDate(0, 0, 1)).
			Scan(&income)

		reports = append(reports, Report{
			Date:        date,
			Income:      income,
			Clicks:      clickCount,
			Impressions: impressionCount,
		})
	}

	c.HTML(http.StatusOK, "reports", gin.H{
		"Reports":     reports,
		"StartDate":   startDateStr,
		"EndDate":     endDateStr,
		"PublisherID": publisherID,
	})
}

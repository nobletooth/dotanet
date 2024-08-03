package publisher

import (
	"bytes"
	"common"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/panel/advertiser"
	"github.com/nobletooth/dotanet/panel/database"
	"gorm.io/gorm"
)

type Publisher struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"column:name;unique;not null"`
	Url         string `gorm:"column:url;unique;not null"`
	Credit      int    `gorm:"column:credit"`
	Clicks      int    `gorm:"column:clicks"`
	Impressions int    `gorm:"column:impressions"`
}

type Report struct {
	Date        time.Time `json:"date"`
	Income      float32   `json:"income"`
	Clicks      int64     `json:"clicks"`
	Impressions int64     `json:"impressions"`
}

func HandlePublisherCreditWithTx(tx *gorm.DB, ad advertiser.Ad, pid int) error {
	creditAddition := int(math.Ceil(ad.Price * 0.2))
	result := tx.Model(&Publisher{}).Where("ID = ?", pid).Update("Credit", gorm.Expr("Credit + ?", creditAddition))
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
	name := c.PostForm("name")
	url := c.PostForm("url")
	publisher := Publisher{Name: name, Credit: 0, Url: url}
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

	script := "<script src=\"" + publisher.Name + ".com/script.js\"></script>"

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

	scriptFilePath := "./publisher/script/template.js"
	scriptContent, err := os.ReadFile(scriptFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read script file"})
		return
	}
	scriptStr := string(scriptContent)
	scriptStr = strings.ReplaceAll(scriptStr, "__PUBLISHER_ID__", idStr)
	scriptStr = strings.ReplaceAll(scriptStr, "__ADSERVER_URL__", *database.AdServerURL)
	modifiedScriptContent := []byte(scriptStr)

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Printf("%v/getadinfo/%v", *database.AdServerURL, idStr)

	c.DataFromReader(http.StatusOK, int64(len(modifiedScriptContent)), "application/javascript", bytes.NewReader(modifiedScriptContent), map[string]string{})
}

func GetPublisherReports(c *gin.Context) {
	publisherIDStr := c.Param("id")
	publisherID, err := strconv.Atoi(publisherIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publisher ID"})
		return
	}

	now := time.Now()
	endDate := now.Truncate(time.Minute)
	startDate := endDate.Add(-1 * time.Hour)

	var reports []Report
	for date := startDate; !date.After(endDate); date = date.Add(time.Minute) {
		var clickCount, impressionCount int64
		var income float32 = 0

		database.DB.Model(&common.ClickedEvent{}).
			Where("pid = ? AND time BETWEEN ? AND ?", publisherID, date, date.Add(time.Minute)).
			Count(&clickCount)

		database.DB.Model(&common.ViewedEvent{}).
			Where("pid = ? AND time BETWEEN ? AND ?", publisherID, date, date.Add(time.Minute)).
			Count(&impressionCount)

		//	query := `
		//SELECT COALESCE(SUM(ads.price * 0.2), 0)
		//FROM "clicked_events"
		//JOIN ads ON clicked_events.ad_id = ads.id
		//WHERE clicked_events.pid = $1
		//AND clicked_events.time BETWEEN $2 AND $3;`
		//
		//	err := database.DB.Raw(query, publisherID, startDate, endDate).Scan(&income).Error
		//	if err != nil {
		//		fmt.Printf("Error executing query: %v", err)
		//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		//		return
		//	}
		database.DB.Table("clicked_events").
			Select("SUM(ads.price)").
			Where("pid = ? AND time BETWEEN ? AND ?", publisherID, date, date.Add(time.Minute)).
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
		"StartDate":   startDate.Format("2006-01-02 15:04:05"),
		"EndDate":     endDate.Format("2006-01-02 15:04:05"),
		"PublisherID": publisherID,
	})
}

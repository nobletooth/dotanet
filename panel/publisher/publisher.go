package publisher

import (
	"bytes"
	"github.com/nobletooth/dotanet/panel/database"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Publisher struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	URL         string `gorm:"column:url;unique;not null"` // Removed space after "url"
	Credit      int    `gorm:"column:credit"`
	Clicks      int    `gorm:"column:clicks"`
	Impressions int    `gorm:"column:impressions"`
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

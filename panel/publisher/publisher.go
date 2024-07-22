package publisher

import (
	"example.com/dotanet/panel/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Publisher struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"unique;not null"`
	Credit      int    `gorm:"column:credit"`
	Script      string `gorm:"column:script"`
	Clicks      int    `gorm:"column:clicks"`
	Impressions int    `gorm:"column:impressions"`
}

func init() {
	err := common.DB.AutoMigrate(&Publisher{})
	if err != nil {
		panic(err)
	}
}

func ListPublishers(c *gin.Context) {
	var publishers []Publisher
	if err := common.DB.Find(&publishers).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "publishers.html", gin.H{"error": "Failed to load publishers"})
		return
	}

	c.HTML(http.StatusOK, "publishers.html", gin.H{
		"Publishers": publishers,
	})
}

func NewPublisherForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create_publisher.html", nil)
}

func CreatePublisherHandler(c *gin.Context) {
	name := c.PostForm("name")
	script := c.PostForm("script")

	publisher := Publisher{Name: name, Credit: 0, Script: script}
	if err := common.DB.Create(&publisher).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": "Failed to create publisher"})
		return
	}

	c.Redirect(http.StatusFound, "/publishers")
}

func ViewPublisherHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Invalid publisher ID"})
		return
	}

	var publisher Publisher
	if err := common.DB.First(&publisher, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{"error": "Publisher not found"})
		return
	}
	c.HTML(http.StatusOK, "view.html", gin.H{
		"Publisher": publisher,
	})
}

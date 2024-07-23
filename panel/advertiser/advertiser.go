package advertiser

import (
	"net/http"
	"strconv"

	"github.com/nobletooth/dotanet/tree/main/panel/common"

	"github.com/gin-gonic/gin"
)

type Entity struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Name   string `gorm:"unique;not null"`
	Credit int    `gorm:"column:credit"`
	Ads    []Ad   `gorm:"foreignKey:AdvertiserId"`
}

type Service interface {
	GetCreditOfAdvertiser(adId int) (Entity, error)
	CreateAdvertiserEntity(name string, credit int)
	ListAllAdvertiserEntities() []Entity
	FindAdvertiserByName(name string) (Entity, error)
	ListAdsByAdvertiser(advertiserId uint) ([]Ad, error)
}

func GetCreditOfAdvertiser(adId int) (Entity, error) {
	var entity Entity

	result := common.DB.First(&entity, adId)
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
	common.DB.Create(&entity)
}

func ListAllAdvertiserEntities() []Entity {
	var advertisers []Entity
	common.DB.Find(&advertisers)
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
	result := common.DB.Where("name = ?", name).First(&entity)
	if result.Error != nil {
		return Entity{}, result.Error
	}
	return entity, nil
}

func ListAdsByAdvertiser(advertiserId uint) ([]Ad, error) {
	var ads []Ad
	result := common.DB.Where("advertiser_id = ?", advertiserId).Find(&ads)
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
	result := common.DB.First(&ad, adID)
	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "index", gin.H{"error": result.Error.Error()})
		return
	}

	c.HTML(http.StatusOK, "edit_ad", gin.H{"Ad": ad})
}

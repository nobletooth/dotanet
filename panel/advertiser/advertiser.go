package advertiser

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

type Entity struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Name   string `gorm:"unique;not null"`
	Credit int    `gorm:"column:credit"`
}

type Service interface {
	GetCreditOfAdvertiser(adId int) (int, error)
	CreateAdvertiserEntity(name string, credit int)
	ListAllAdvertiserEntities() []Entity
	FindAdvertiserByName(name string) (Entity, error)
}

func GetCreditOfAdvertiser(adId int) (int, error) {
	var entity Entity
	result := DB.First(&entity, adId)
	if result.Error != nil {
		return 0, result.Error
	}
	return entity.Credit, nil
}

func CreateAdvertiserEntity(name string, credit int) {
	entity := Entity{
		Name:   name,
		Credit: credit,
	}
	DB.Create(&entity)
}

func ListAllAdvertiserEntities() []Entity {
	var advertisers []Entity
	DB.Find(&advertisers)
	return advertisers
}

func FindAdvertiserByName(name string) (Entity, error) {
	var entity Entity
	result := DB.Where("name = ?", name).First(&entity)
	if result.Error != nil {
		return Entity{}, result.Error
	}
	return entity, nil
}

func LoadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", templatesDir+"/index.html")
	r.AddFromFiles("create_advertiser", templatesDir+"/create_advertiser.html")
	r.AddFromFiles("advertiser_credit", templatesDir+"/advertiser_credit.html")
	return r
}

func ListAdvertisers(c *gin.Context) {
	advertisers := ListAllAdvertiserEntities()
	c.HTML(http.StatusOK, "index", gin.H{"Advertisers": advertisers})
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

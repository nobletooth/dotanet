package advertiser

import "example.com/dotanet/db"

type AdvertiserRequestDto struct {
 Creadit int    json:"amount"
 Name    string json:"name"
}

type AdvertiserEntity struct {
 ID     uint   gorm:"primaryKey;autoIncrement" // Unique identifier, auto-incrementing
 Name   string gorm:"unique;not null"          // Unique name, optional but ensures that the name is unique
 Credit int    gorm:"column:credit"            // Fixed typo for the column name
}

type Advertiserservice interface {
 GetCreaditOfAdvertiser(adId int) (int, error)
 CreateAdvertiserEntity(AdvertiserRequestDto)
 ListAllAdvertiserEntity() []AdvertiserEntity
 FindAdvertiserByName(name string) (AdvertiserEntity, error)
}

type AdvertiService struct {
 db *db.Database
}

func NewAdvertiserService(db *db.Database) AdvertiService {
 return AdvertiService{db: db}
}

func (p *AdvertiService) CreateAdvertiserEntity(dto AdvertiserRequestDto) {
 entity := AdvertiserEntity{
  Credit: dto.Creadit,
  Name:    dto.Name,
 }
 p.db.DB.Save(&entity)
}

func (p *AdvertiService) ListAllAdvertiserEntity() []AdvertiserEntity {
 var advertisers []AdvertiserEntity
 p.db.DB.Find(&advertisers)
 return advertisers
}

func (p *AdvertiService) GetCreaditOfAdvertiser(adId int) (int, error) {
 var entity AdvertiserEntity
 result := p.db.DB.First(&entity, adId)
 if result.Error != nil {
  return 0, result.Error
 }
 return entity.Creadit, nil
}

func (p *AdvertiService) FindAdvertiserByName(name string) (AdvertiserEntity, error) {
 var entity AdvertiserEntity
 result := p.db.DB.Where("name = ?", name).First(&entity)
 if result.Error != nil {
  return AdvertiserEntity{}, result.Error
 }
 return entity, nil
}

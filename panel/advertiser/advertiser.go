package advertiser

import "example.com/dotanet/db"

type AdvertiserRequestDto struct {
	ID      int `json:"id"`
	Creadit int `json:"amount"`
}

type AdvertiserEntity struct {
	ID      int  `gorm:"column:id;primary_key"`
	Creadit bool `gorm:"column:creadit"`
}

type Advertiserservice interface {
	// AddTask(dto.TodoRequestBody) (entitys.TodoLists, error)
	// ListTasks() ([]entitys.TodoLists, error)
	// UpdateTask(int, dto.TodoRequestBody) error
	// RemoveTask(int) error
	// GetTask(int) (entitys.TodoLists, error)
	CreateAdvertiserEntity(AdvertiserEntity) AdvertiserEntity
}

type AdvertiService struct {
	db *db.Database
}

func NewAdvertiserService(db *db.Database) AdvertiService {
	AdvertiSerservice := AdvertiService{db: db}
	return AdvertiSerservice
}

func (p *AdvertiService) CreateAdvertiserEntity(entity AdvertiserEntity) {
	p.db.DB.Create(&entity)
}

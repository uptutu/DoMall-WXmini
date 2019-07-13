package Brand

import (
	"do-mall/models"
	"do-mall/pkg/logging"
)

type Brand struct {
	models.Model
	Name string `json:"name"`
	Logo string `json:"logo"`
}

func (Brand) TableName() string {
	return "brands"
}

func Create(brand *Brand) bool {
	if models.DB.NewRecord(brand) {
		models.DB.Create(brand)
		return !models.DB.NewRecord(brand)
	}
	
	return false
}

func Update(brand *Brand, data interface{}) bool {
	if err := models.DB.Debug().Model(brand).Update(data).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func Index(pageNum int, pageSize int, maps interface {}) (brands []Brand) {
	models.DB.Where(maps).Offset(pageNum).Limit(pageSize).Find(&brands)
	return
}

func Total(maps interface{})(count int){
	models.DB.Model(&Brand{}).Where(maps).Count(&count)
	return
}
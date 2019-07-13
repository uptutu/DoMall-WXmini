package Product

import (
	"do-mall/models"
	"do-mall/pkg/logging"
)

type Product struct {
	models.Model

	Title        string  `json:"title"`
	Cover        string  `json:"cover"`
	Carousel     string  `json:"carousel"`
	Brand        string  `json:"brand"`
	Series       string  `json:"series"`
	Price        float32 `json:"price"`
	SellingPrice float32 `json:"selling_price"`
	Cost         float32 `json:"cost"`
	Tags         string  `json:"tags"`
	Sales        int     `json:"sales"`
	Inventory    int     `json:"inventory"`
	Status       int     `json:"status"`
	OnSale       int     `json:"on_sale"`
}

func TableName() string {
	return "products"
}

func Create(product *Product) bool {
	if models.DB.NewRecord(product) {
		models.DB.Create(product)
		return models.DB.NewRecord(product)
	}

	return false
}

func Update(product *Product, data interface{}) bool {
	if err := models.DB.Debug().Model(product).Update(data).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func Show(id int)(product Product) {
	if err := models.DB.Debug().Where(map[string]interface{}{"id":id}).First(&product).Error; err != nil {
		logging.Info(err)
	}
	return
}

func List(pageNum int, pageSize int, maps interface{}) (products []Product) {
	models.DB.Where(maps).Offset(pageNum).Limit(pageSize).Find(&products)
	return
}

func SearchInTitle(pageNum int, pageSize int, data interface{})(products []Product) {
	models.DB.Where("title LIKE ?", data).Offset(pageNum).Limit(pageSize).Find(&products)
	return
}

func SearchInTag(pageNum int, pageSize int, data interface{})(products []Product) {
	models.DB.Where("Tag LIKE ?", data).Offset(pageNum).Limit(pageSize).Find(&products)
	return
}

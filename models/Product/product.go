package Product

import (
	"do-mall/models"
	"do-mall/pkg/logging"
	"fmt"
)

type Product struct {
	models.Model

	Title        string  `json:"title"`
	Cover        string  `json:"cover"`
	Carousel     string  `json:"carousel"`
	Brand        string  `json:"brand"`
	Series       string  `json:"series"`
	Price        float64 `json:"price"`
	SellingPrice float64 `json:"selling_price"`
	Cost         float64 `json:"cost"`
	Tags         string  `json:"tags"`
	Sales        int     `json:"sales"`
	Inventory    int     `json:"inventory"`
	Status       int     `json:"status"`
	OnSale       int     `json:"on_sale"`
}

func (Product)TableName() string {
	return "products"
}

func Create(product *Product) bool {
	if models.DB.NewRecord(product) {
		models.DB.Create(product)
		return !models.DB.NewRecord(product)
	}

	return false
}

func Update(data *Product) bool {
	if data.ID == 0 {
		return false
	}
	product := Product{}
	models.DB.Where("id = ?", data.ID).First(&product)
	if err := models.DB.Debug().Model(product).Update(data).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func Destroy(id int) bool {
	if err := models.DB.Debug().Delete(Product{}, "id = ?", id).Error; err != nil {
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

func List(pageNum int, pageSize int, maps interface{}) (products []Product, count int) {
	models.DB.Where(maps).Offset(pageNum).Limit(pageSize).Find(&products)
	models.DB.Model(Product{}).Where(maps).Count(&count)

	return
}

func SearchInTitle(pageNum int, pageSize int, data string)(products []Product, count int) {
	models.DB.Debug().Where("title LIKE ?", data).Offset(pageNum).Limit(pageSize).Find(&products)
	models.DB.Debug().Model(Product{}).Where("`title` LIKE ?", data).Count(&count)
	return
}

func SearchInTag(pageNum int, pageSize int, data string)(products []Product, count int) {
	data = fmt.Sprintf("%%%s%%", data)
	models.DB.Debug().Where("tags LIKE ?", data).Offset(pageNum).Limit(pageSize).Find(&products)
	models.DB.Debug().Model(Product{}).Where("`tags` LIKE ?", data).Count(&count)
	return
}



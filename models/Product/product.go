package Product

import (
	"do-mall/models"
	"do-mall/pkg/logging"
	"fmt"
	"time"
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

type Inventory struct {
	PId       int        `json:"p_id"`
	F         int        `json:"f"`
	Xxs       int        `json:"xxs"`
	Xs        int        `json:"xs"`
	S         int        `json:"s"`
	M         int        `json:"m"`
	L         int        `json:"l"`
	Xl        int        `json:"xl"`
	Xxl       int        `json:"xxl"`
	S35       int        `json:"s_35"`
	S36       int        `json:"s_36"`
	S37       int        `json:"s_37"`
	S38       int        `json:"s_38"`
	S39       int        `json:"s_39"`
	S40       int        `json:"s_40"`
	S405      int        `json:"s_405"`
	S41       int        `json:"s_41"`
	S415      int        `json:"s_415"`
	S42       int        `json:"s_42"`
	S425      int        `json:"s_425"`
	S43       int        `json:"s_43"`
	S435      int        `json:"s_435"`
	S44       int        `json:"s_44"`
	S445      int        `json:"s_445"`
	S45       int        `json:"s_45"`
	S46       int        `json:"s_46"`
	S47       int        `json:"s_47"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Product) TableName() string {
	return "products"
}

func (Inventory) TableName() string {
	return "inventories"
}

func Create(product *Product) bool {
	if models.DB.NewRecord(product) {
		models.DB.Create(product)
		if !models.DB.NewRecord(product) {
			inv := Inventory{}
			if err := models.DB.Debug().Create(inv).Error; err != nil {
				logging.Info(err)
				return false
			} else {
				return true
			}
		}
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
	if !DestroyInventory(id) {
		return false
	}
	if err := models.DB.Debug().Delete(Product{}, "id = ?", id).Error; err != nil {
		logging.Info(err)
		return false
	}

	return true
}

func Show(id int) (product Product) {
	if err := models.DB.Debug().Where(map[string]interface{}{"id": id}).First(&product).Error; err != nil {
		logging.Info(err)
	}
	return
}

func List(pageNum int, pageSize int, maps interface{}) (products []Product, count int) {
	models.DB.Where(maps).Offset(pageNum).Limit(pageSize).Find(&products)
	models.DB.Model(Product{}).Where(maps).Count(&count)

	return
}

func SearchInTitle(pageNum int, pageSize int, data string) (products []Product, count int) {
	models.DB.Debug().Where("title LIKE ?", data).Offset(pageNum).Limit(pageSize).Find(&products)
	models.DB.Debug().Model(Product{}).Where("`title` LIKE ?", data).Count(&count)
	return
}

func SearchInTag(pageNum int, pageSize int, data string) (products []Product, count int) {
	data = fmt.Sprintf("%%%s%%", data)
	models.DB.Debug().Where("tags LIKE ?", data).Offset(pageNum).Limit(pageSize).Find(&products)
	models.DB.Debug().Model(Product{}).Where("`tags` LIKE ?", data).Count(&count)
	return
}

func UpdateInventory(pId int, data *Inventory) bool {
	inventory := GetInventory(pId)
	if err := models.DB.Debug().Model(inventory).Update(*data).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func GetInventory(pId int) (inventory Inventory) {
	models.DB.Debug().Where("p_id = ?", pId).First(&inventory)
	return
}

func DestroyInventory(pid int) bool {
	if err := models.DB.Debug().Delete(Inventory{}, "id = ?", pid).Error; err != nil {
		logging.Info(err)
		return false
	}

	return true
}

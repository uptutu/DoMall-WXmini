package User

import (
	"do-mall/models"
	"do-mall/models/Product"
	"do-mall/pkg/logging"
	"time"
)

type Cart struct {
	ID     int    `gorm:"primary_key" json:"id"`
	UserId int    `json:"user_id"`
	PId    int    `json:"p_id"`
	OId    int    `json:"o_id"`
	Number int    `json:"number"`
	Size   string `json:"size"`
	Status int    `json:"status"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Cart) TableName() string {
	return "carts"
}

func PutInCart(userId, pId int, size string) bool {
	if IsInCart(userId, pId, size) {
		cart := QueryCart(userId, pId)
		NumberIncrease(cart.ID)
		return true
	} else {
		if ok := CreateCartRow(userId, pId, size); ok {
			return true
		}
		return false
	}
}

func DropFromCart(cId int) bool {
	cart := Cart{ID: cId}
	if err := models.DB.Debug().Delete(&cart, "id = ?", cart.ID).Error; err != nil {
		logging.Info(err)
		return false
	}

	return true
}

func CartsProducts(userId int) (products []Product.Product) {
	db := models.DB.Table("carts").Select("p_id").Where("user_id = ?", userId)
	rows, err := db.Rows()
	if err != nil {
		logging.Info(err)
	}
	var pids []int
	for rows.Next() {
		var pid int
		_ = rows.Scan(&pid)
		pids = append(pids, pid)
	}
	models.DB.Where("id in (?)", pids).Find(&products)
	return
}

func CreateCartRow(userId, pId int, size string) bool {
	cart := Cart{UserId: userId, PId: pId, Number: 1, Size: size, CreatedAt: time.Now()}
	if err := models.DB.Debug().Create(&cart).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func QueryCart(userId, pId int) (cart Cart) {
	models.DB.Debug().Model(Cart{}).Where("user_id = ? AND p_id = ?", userId, pId).First(&cart)
	return
}

func QueryCartById(cId int) (cart Cart) {
	models.DB.Debug().Model(Cart{}).Where("id = ?", cId).First(&cart)
	return
}

func IsInCart(userId, pId int, size string) bool {
	cars := ListCart(userId)
	for _, item := range cars {
		if item.PId == pId && item.Size == size {
			return true
		}
	}
	return false
}

func ListCart(userId int) (carts []Cart) {
	models.DB.Debug().Model(Cart{}).Where("user_id = ?", userId).Find(&carts)
	return
}

func NumberIncrease(cId int) bool {
	cart := QueryCartById(cId)
	cart.Number += 1
	if err := models.DB.Debug().Model(cart).UpdateColumn("number", cart.Number).Error; err != nil {
		return false
	}
	return true
}

func NumberDecrease(cId int) bool {
	cart := QueryCartById(cId)
	cart.Number -= 1
	if cart.Number <= 0 {
		DropFromCart(cId)
	}
	if err := models.DB.Debug().Model(cart).UpdateColumn("number", cart.Number).Error; err != nil {
		return false
	}

	return true
}

func OrderCreated(ids interface{}, oid int) bool {
	if err := models.DB.Debug().Model(Cart{}).Where("id in (?)", ids).UpdateColumn("o_id", oid).Error; err != nil {
		return false
	}

	return true
}

func QueryCarts(ids interface{}) (carts []Cart) {
	models.DB.Debug().Model(Cart{}).Where("id in (?)", ids).Find(&carts)
	return
}

func QueryCartByOid(oid int) (carts []Cart) {
	models.DB.Debug().Model(Cart{}).Where("o_id = ?",oid).Find(&carts)
	return
}

package Order

import (
	"do-mall/models"
	"do-mall/models/Wallet"
	"do-mall/pkg/logging"
	"time"
)

type Order struct {
	models.Model
	UserId       int       `json:"user_id"`
	ProvinceName string    `json:"province_name"`
	CityName     string    `json:"city_name"`
	CountyName   string    `json:"county_name"`
	DetailInfo   string    `json:"detail_info"`
	PostalCode   string    `json:"postal_code"`
	UserName     string    `json:"user_name"`
	TelNumber    string    `json:"tel_number"`
	ExpressTitle string    `json:"express_title"`
	ExpressCode  string    `json:"express_code"`
	ExpressNo    string    `json:"express_no"`
	ExpressTime  time.Time `json:"express_time"`
	Total        float64   `json:"total"`
	SumPay       float64   `json:"sum_pay"`
	Status       int       `json:"status"`
	PushedAt     time.Time `json:"pushed_at"`
}

func (Order) TableName() string {
	return "orders"
}

func Created(order *Order) bool {
	if err := models.DB.Debug().Create(order).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func Settlement(userId, oid int) bool {
	order := QueryOrderById(oid)
	if Wallet.CostBalance(order.SumPay, userId) {
		if err := models.DB.Debug().Model(order).UpdateColumn("status", 1).Error; err != nil {
			logging.Info(err)
			return false
		}
		return true
	} else {
		return false
	}
}

func Update(id int, order *Order) bool {
	if err := models.DB.Debug().Model(Order{}).Where("id = ?", id).Update(*order).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func Destroy(oid int) bool {
	if err := models.DB.Unscoped().Delete(Order{}, "id = ?", oid).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func Done(oid int) bool {
	if err := models.DB.Delete(Order{}, "id = ?", oid).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func QueryOrderByUserId(userId int) (order []Order) {
	models.DB.Debug().Model(Order{}).Where("user_id = ?", userId).Find(&order)
	return
}

func QueryOrderById(id int) (order Order) {
	models.DB.Debug().Model(Order{}).Where("id = ?", id).First(&order)
	return
}

func IsOwner(userId, id int) bool {
	var data int
	models.DB.Debug().Model(Order{}).Select("1").Where("id = ? AND user_id = ?", id, userId).First(&data)
	if data != 0 {
		return true
	}
	return false
}
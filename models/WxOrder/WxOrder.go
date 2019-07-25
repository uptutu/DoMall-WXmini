package WxOrder

import (
	"do-mall/models"
	"do-mall/pkg/logging"
)

type WxOrder struct {
	models.Model
	UserId        int     `json:"user_id"`
	AppID         int     `json:"app_id"`
	Body          string  `json:"body"`
	OutTradeNo    string  `json:"out_trade_no"`
	NonceStr      string  `json:"nonce_str"`
	Sign          string  `json:"sign"`
	SumPay        float64 `json:"sum_pay"`
	TotalFee      int     `json:"total_fee"`
	Detail        string  `json:"detail"`
	Attach        string  `json:"attach"`
	TransactionId string  `json:"transaction_id"`
}

func (WxOrder) TableName() string {
	return "wxorders"
}

func Create(order *WxOrder) bool {
	if models.DB.NewRecord(*order) {
		models.DB.Debug().Create(*order)
		return !models.DB.NewRecord(*order)
	}

	return false
}

func Update(order *WxOrder) bool {
	if err := models.DB.Debug().Model(WxOrder{}).Update(*order).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func QueryByTradeNo(no string) (order WxOrder) {
	models.DB.Debug().Model(&order).Where("tradeNo = ?", no).First(&order)
	return
}

func Destroy(order *WxOrder) bool {
	if err := models.DB.Debug().Model(WxOrder{}).Delete(*order).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}
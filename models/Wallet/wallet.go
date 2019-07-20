package Wallet

import (
	"do-mall/models"
	"do-mall/models/User"
	"do-mall/pkg/logging"
)

type Wallet struct {
	models.Model
	Balance float64 `json:"balance"`
	Coin    float64 `json:"coin"`
}

func (Wallet) TableName() string {
	return "users"
}

func TopUpBalance(amount float64, userId int) bool {
	user := User.GetInfo(userId)
	wallet := Wallet{}
	wallet.Balance = user.Balance + amount
	if err := models.DB.Debug().Model(wallet).UpdateColumn("balance", wallet.Balance).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func CostBalance(amount float64, userId int)  bool {
	user := User.GetInfo(userId)
	wallet := Wallet{}
	wallet.Balance = user.Balance - amount
	if wallet.Balance < 0 {
		return false
	}
	if err := models.DB.Debug().Model(wallet).UpdateColumn("balance", wallet.Balance).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

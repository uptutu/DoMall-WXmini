package User

import (
	"do-mall/models"
	"do-mall/pkg/logging"
)

type Address struct {
	ID           int    `json:"id"`
	UserId       int    `json:"user_id"`
	ProvinceName string `json:"province_name"`
	CityName     string `json:"city_name"`
	CountyName   string `json:"county_name"`
	DetailInfo   string `json:"detail_info"`
	PostalCode   string `json:"postal_code"`
	NationalCode string `json:"national_code"`
	UserName     string `json:"user_name"`
	TelNumber    string `json:"tel_number"`
	Default      int    `json:"default"`
}

func (Address) TableName() string {
	return "addresses"
}

func CreateAddress(address *Address) bool {
	if err := models.DB.Debug().Create(address).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func DestroyAddress(id int) bool {
	if err := models.DB.Delete(Address{}, "id = ?", id).Error;err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func SelectAddress(id int) (address Address) {
	models.DB.Debug().Model(Address{}).Where("id = ?", id).First(&address)
	return
}

func AddressesOfUser(userId int) (addresses []Address) {
	if err := models.DB.Debug().Model(Address{}).Where("user_ud = ?", userId).Find(&addresses).Error; err != nil {
		logging.Info(err)
	}
	return
}

func SetAddressDefault(id int) bool {
	address := SelectAddress(id)
	SetAllAddressNotDefault(address.UserId)
	if err := models.DB.Debug().Model(address).UpdateColumn("default", 1).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func SetAllAddressNotDefault(userId int) bool {
	if err := models.DB.Debug().Model(Address{}).Where("user_id = ?", userId).UpdateColumn("default", 0).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

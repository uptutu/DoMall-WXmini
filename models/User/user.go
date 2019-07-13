package User

import (
	"do-mall/models"
	"do-mall/pkg/logging"
)

type User struct {
	models.Model

	Password     string  `json:"-"`
	Unionid      string  `json:"-"`
	Openid       string  `json:"-"`
	Nickname     string  `json:"nickname"`
	Avatar       string  `json:"avatar"`
	Sex          int     `json:"sex"`
	Mobile       string  `json:"mobile"`
	Introduction string  `json:"introduction"`
	Balance      float32 `json:"balance"`
	Coin         float32 `json:"coin"`
}

func (User) TableName() string {
	return "users"
}

func GetInfo(id int) (user User) {
	if err := models.DB.Debug().Where(map[string]interface{}{"id": id}).Find(&user).Error; err != nil {
		logging.Info(err)
	}
	return
}

func Login(user User) bool {
	var find User
	if models.DB.Where(user).Select("id").First(&find); find.ID > 0 {
		return true
	}
	return false
}

func CreateByPasswd(user *User) bool {
	if models.DB.NewRecord(user) {
		models.DB.Create(user)
		return !models.DB.NewRecord(user)
	}

	return false
}

func Update(user *User, data interface{}) bool {
	if err := models.DB.Debug().Model(user).Update(data).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

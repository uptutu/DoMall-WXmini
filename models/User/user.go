package User

import "do-mall/models"

type User struct {
	models.Model

	Password     string  `json:"password"`
	Unionid      string  `json:"unionid"`
	Openid       string  `json:"openid"`
	Nickname     string  `json:"nickname"`
	Avatar       string  `json:"avatar"`
	Sex          int     `json:"sex"`
	Mobile       string  `json:"mobile"`
	Introduction string  `json:"introduction"`
	Balance      float32 `json:"balance"`
	Coin         float32 `json:"coin"`

}

func GetInfo(id int) (user User) {
	models.DB.Where("id = ?", id).First(&user)
	return
}

func CreateByPasswd(user *User) bool {
	if models.DB.NewRecord(user) {
		models.DB.Create(user)
		return !models.DB.NewRecord(user)
	}

	return false

}

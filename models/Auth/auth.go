package Auth

import "do-mall/models"

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Unionid  string `json:"-"`
	Openid   string `json:"-"`
}

type AuthAdmin struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAndReturnId(mobile, password string) int {
	var auth Auth
	models.DB.Select("id").Where(Auth{Mobile: mobile, Password: password}).First(&auth)
	return int(auth.ID)
}

func (Auth) TableName() string {
	return "users"
}

func (AuthAdmin) TableName() string {
	return "admins"
}

func CheckAdmin(username, password string) int {
	var auth AuthAdmin
	models.DB.Debug().Select("id").Where(AuthAdmin{Username: username, Password: password}).First(&auth)
	return int(auth.ID)
}

func CheckOpenid(openid string) int {
	var auth Auth
	models.DB.Debug().Select("id").Where(Auth{Openid:openid}).First(&auth)
	return int(auth.ID)
}

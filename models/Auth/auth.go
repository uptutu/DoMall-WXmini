package Auth

import "do-mall/models"

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

func CheckAuth(mobile, password string) bool {
	var auth Auth
	models.DB.Select("id").Where(Auth{Mobile: mobile, Password: password}).First(&auth)
	if auth.ID > 0 {
		return true
	}

	return false
}

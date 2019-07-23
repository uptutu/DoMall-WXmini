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
	SessionKey   string  `json:"-"`
	Nickname     string  `json:"nickname"`
	Avatar       string  `json:"avatar"`
	Sex          int     `json:"sex"`
	Mobile       string  `json:"mobile"`
	Introduction string  `json:"introduction"`
	Balance      float64 `json:"balance"`
	Coin         float64 `json:"coin"`
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
		models.DB.Debug().Create(user)
		return !models.DB.NewRecord(user)
	}

	return false
}

func Update(user *User, data interface{}) bool {
	if err := models.DB.Debug().Model(user).Where("id = ?", user.ID).Update(data).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func UpdateColumn(user *User,column string, data interface{}) bool {
	if err := models.DB.Debug().Model(User{}).Where("id = ?", user.ID).UpdateColumn(column, data).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func Create(user *User) bool {
	if models.DB.NewRecord(*user) {
		models.DB.Debug().Create(*user)
		return !models.DB.NewRecord(*user)
	}

	return false
}

func AddUnionid(userId int, unionid string) bool {
	if err := models.DB.Debug().Model(User{}).Where("id = ?", userId).UpdateColumn("unionid", unionid).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func QueryUserByUnionid(unionid string) (user User) {
	models.DB.Debug().Model(User{}).Where("unionid = ?", unionid).First(&user)
	return
}

func QueryUserByOpenid(openid string) (user User) {
	models.DB.Debug().Model(User{}).Where("openid = ?", openid).First(&user)
	return
}

func QueryUserByid(id int) (user User) {
	models.DB.Debug().Model(User{}).Where("id = ?", id).First(&user)
	return
}

func PutSsk(openid, ssk string) bool {
	user := QueryUserByOpenid(openid)
	user.SessionKey = ssk
	if err := models.DB.Debug().Model(user).Update(user).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

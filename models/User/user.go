package User

import (
	"do-mall/models"
	"do-mall/pkg/logging"
	"time"
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
		models.DB.Debug().Create(user)
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

// ------------------------------------------
// 用户收藏
// ------------------------------------------

type Favorite struct {
	PId   int
	UserID  int
	CreatedAt time.Time
}

func (Favorite) TableName() string {
	return "favorites"
}

func AddFavorite(userId, pId int) bool {
	favs := ListFavorite(userId)
	for _, item:= range favs {
		if item.PId == pId {
			return false
		}
	}
	fav := Favorite{UserID: userId, PId: pId, CreatedAt:time.Now()}
	if models.DB.NewRecord(fav) {
		models.DB.Debug().Create(fav)
		return !models.DB.NewRecord(fav)
	}
	return false
}

func DestroyFavorite(userId, pId int) bool {
	fav := Favorite{UserID:userId, PId:pId}
	if err := models.DB.Debug().Unscoped().Model(&fav).Delete(&fav).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func ListFavorite(userId int) (favorites []Favorite) {
	models.DB.Model(Favorite{}).Where("user_id = ?", userId).Find(favorites)
	return
}
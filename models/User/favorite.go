package User

import (
	"do-mall/models"
	"do-mall/models/Product"
	"do-mall/pkg/logging"
	"time"
)

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
	if err := models.DB.Debug().Create(fav).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func DestroyFavorite(userId, pId int) bool {
	fav := Favorite{UserID:userId, PId:pId}
	if err := models.DB.Debug().Unscoped().Delete(&fav, "user_id = ? AND p_id = ?", fav.UserID, fav.PId).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

func ListFavorite(userId int) (favorites []Favorite) {
	models.DB.Model(Favorite{}).Where("user_id = ?", userId).Find(&favorites)
	return
}

func ShowFavorites(pageNum int, pageSize int, userId int) (Products []Product.Product, count int) {
	db := models.DB.Table("favorites").Select("p_id").Where("user_id = ?", userId)
	rows, err := db.Rows()
	if err != nil {
		logging.Info(err)
	}
	var pids []int
	for rows.Next(){
		var pid int
		_ = rows.Scan(&pid)
		pids = append(pids, pid)
	}
	var product Product.Product
	models.DB.Where("id in (?)", pids).Offset(pageNum).Limit(pageSize).Find(&Products)
	models.DB.Model(&product).Where("id in (?)", pids).Count(&count)
	return
}

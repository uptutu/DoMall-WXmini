package ProductController

import (
	"do-mall/models"
	"do-mall/models/Product"
	"do-mall/pkg/e"
	"do-mall/pkg/setting"
	"do-mall/pkg/util"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})
	var msg string

	condition := struct {
		Where map[string]interface{}
		Like  map[string]interface{}
	}{}

	if tag := c.Query("tag"); tag != "" {
		condition.Like["tag"] = tag
	}

	if brand := c.Query("brand"); brand != "" {
		condition.Where["brand"] = brand
	}

	if series := c.Query("series"); series != "" {
		condition.Where["series"] = series
	}

	if title := c.Query("title"); title != "" {
		condition.Like["title"] = title
	}


	data["lists"] = Product.List(util.GetPage(c), setting.PageSize, condition.Where)

}

func Show(c *gin.Context) {
	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})
	var msg string

	id := com.StrTo(c.Param("id")).MustInt()
	if product := Product.Show(id); !(product.ID >= 0) {
		code = e.BAD_REQUEST
		msg = "请求不到指定id数据"
	} else {
		code = e.OK
		msg = e.GetMsg(code)
		data["product"] = product
	}

	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})

}

package ProductController

import (
	"do-mall/models/Product"
	"do-mall/pkg/e"
	"do-mall/pkg/setting"
	"do-mall/pkg/util"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
)

//
// @Summary 商品列表
// @Produce json
// @param string query false "brand"
// @param state query string false "series"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"Ok"}"
// @Failure 400 {string} json "{"code":400,"msg":"Bad Request","data":{}}"
// @Router /api/v1/product/index [get]
func Index(c *gin.Context) {

	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string

	where := make(map[string]interface{})

	if brand := c.Query("brand"); brand != "" {
		where["brand"] = brand
	}
	if series := c.Query("series"); series != "" {
		where["series"] = series
	}

	data["lists"], data["total"] = Product.List(util.GetPage(c), setting.PageSize, where)
	if _, ok := data["lists"]; ok {
		code = e.OK
	}
	msg = e.GetMsg(code)

	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})

}

func Search(c *gin.Context) {
	code := e.BAD_REQUEST
	var msg string
	data := make(map[string]interface{})

	if tag := c.Query("tag"); tag != "" {
		data["lists"], data["total"] = Product.SearchInTag(util.GetPage(c), setting.PageSize, tag)

	} else if title := c.Query("title"); title != "" {
		data["lists"], data["total"] = Product.SearchInTitle(util.GetPage(c), setting.PageSize, tag)
	}

	if _, ok := data["lists"]; ok {
		code = e.OK
	}
	msg = e.GetMsg(code)

	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})

}

func Show(c *gin.Context) {
	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})
	var msg string

	id := com.StrTo(c.Param("id")).MustInt()
	if product := Product.Show(id); !(product.ID > 0) {
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

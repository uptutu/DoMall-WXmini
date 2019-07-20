package CartController

import (
	"do-mall/models/Product"
	"do-mall/models/User"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"do-mall/pkg/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Index(c *gin.Context) {
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID

	data["products"] = User.CartsProducts(userId)
	code := e.OK
	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func Create(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}

	valid.Required(c.PostForm("pId"), "pId").Message("pId 必须")
	valid.Numeric(c.PostForm("pId"), "pId").Message("pId 必须有效")
	valid.Required(c.PostForm("size"), "size").Message("size 必须")
	if valid.HasErrors() {
		errorData := make(map[string]interface{})
		for index, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
			errorData[strconv.Itoa(index)] = map[string]interface{}{err.Key: err.Message}
		}
		data["error"] = errorData
	}

	if _, ok := data["error"]; !ok {
		pId := com.StrTo(c.PostForm("pId")).MustInt()
		size := c.PostForm("size")
		Product.DecreaseInventory(pId, size)
		if User.PutInCart(userId, pId, size) {
			code = e.NO_CONTENT
		}
	} else {
		code = e.INTERNAL_SERVER_ERROR
	}

	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func Destroy(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}

	valid.Required(c.Param("id"), "cId").Message("cart Id 参数值必须")
	valid.Numeric(c.Param("id"), "cId").Message("cart Id 必须是有效数值")

	if valid.HasErrors() {
		code = e.BAD_REQUEST
		errorData := make(map[string]interface{})
		for index, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
			errorData[strconv.Itoa(index)] = map[string]interface{}{err.Key: err.Message}
		}
		data["error"] = errorData
	}

	if _, ok := data["error"]; !ok {
		cId := com.StrTo(c.Param("id")).MustInt()
		cart := User.QueryCartById(cId)
		if cart.UserId != userId {
			code = e.UNAUTHORIZED
		} else if User.DropFromCart(cId) {
			code = e.NO_CONTENT
		}
	} else {
		code = e.INTERNAL_SERVER_ERROR
	}

	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func Decrease(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}

	valid.Required(c.Param("id"), "cId").Message("cart Id 参数值必须")
	valid.Numeric(c.Param("id"), "cId").Message("cart Id 必须是有效数值")

	if valid.HasErrors() {
		code = e.BAD_REQUEST
		errorData := make(map[string]interface{})
		for index, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
			errorData[strconv.Itoa(index)] = map[string]interface{}{err.Key: err.Message}
		}
		data["error"] = errorData
	}

	if _, ok := data["error"]; !ok {
		cId := com.StrTo(c.Param("id")).MustInt()
		cart := User.QueryCartById(cId)
		if cart.UserId != userId {
			code = e.UNAUTHORIZED
		} else if User.NumberDecrease(cId) {
			code = e.NO_CONTENT
		}
	} else {
		code = e.INTERNAL_SERVER_ERROR
	}

	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

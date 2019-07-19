package AddressController

import (
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

	addresses := User.AddressesOfUser(userId)

	code := e.OK
	data["addresses"] = addresses

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

	valid.Required(c.PostForm("provinceName"), "provinceName").Message("必须")
	valid.Required(c.PostForm("cityName"), "cityName").Message("cityName 必须")
	valid.Required(c.PostForm("countyName"), "countyName").Message("countyName 必须")
	valid.Required(c.PostForm("detailInfo"), "detailInfo").Message("detailInfo 必须")
	valid.Required(c.PostForm("postalCode"), "postalCode").Message("postalCode 必须")
	valid.Required(c.PostForm("nationalCode"), "nationalCode").Message("nationalCode 必须")
	valid.Required(c.PostForm("userName"), "userName").Message("userName 必须")
	valid.Required(c.PostForm("telNumber"), "telNumber").Message("telNumber 必须")
	valid.Phone(c.PostForm("telNumber"), "telNumber").Message("telNumber 必须是有效电话号码")

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
		address := User.Address{}
		address.UserId = userId
		address.ProvinceName = c.PostForm("provinceName")
		address.CityName = c.PostForm("cityName")
		address.CountyName = c.PostForm("countyName")
		address.DetailInfo = c.PostForm("detailInfo")
		address.PostalCode = c.PostForm("postalCode")
		address.NationalCode = c.PostForm("nationalCode")
		address.UserName = c.PostForm("userName")
		address.TelNumber = c.PostForm("telNumber")
		if User.CreateAddress(&address) {
			code = e.CREATED
		} else {
			code = e.INTERNAL_SERVER_ERROR
		}

		msg = e.GetMsg(code)
	}

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
	addId := com.StrTo(c.Param("id")).MustInt()

	add := User.SelectAddress(addId)
	if add.UserId != userId {
		code = e.UNAUTHORIZED
		msg = e.GetMsg(code)
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
	}
	if User.DestroyAddress(addId) {
		code = e.NO_CONTENT
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

func SetDefault(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	addId := com.StrTo(c.Param("id")).MustInt()
	if addId <= 0 {
		code = e.BAD_REQUEST
		msg = "请输入有效id"
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
	}
	add := User.SelectAddress(addId)
	if add.UserId != userId {
		code = e.UNAUTHORIZED
		msg = e.GetMsg(code)
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
	}
	if User.SetAddressDefault(addId) {
		code = e.NO_CONTENT
	}
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})

}

func Update(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	AddId := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}

	valid.Required(c.PostForm("provinceName"), "provinceName").Message("必须")
	valid.Required(c.PostForm("cityName"), "cityName").Message("cityName 必须")
	valid.Required(c.PostForm("countyName"), "countyName").Message("countyName 必须")
	valid.Required(c.PostForm("detailInfo"), "detailInfo").Message("detailInfo 必须")
	valid.Required(c.PostForm("postalCode"), "postalCode").Message("postalCode 必须")
	valid.Required(c.PostForm("nationalCode"), "nationalCode").Message("nationalCode 必须")
	valid.Required(c.PostForm("userName"), "userName").Message("userName 必须")
	valid.Required(c.PostForm("telNumber"), "telNumber").Message("telNumber 必须")
	valid.Phone(c.PostForm("telNumber"), "telNumber").Message("telNumber 必须是有效电话号码")


	if valid.HasErrors() {
		code = e.BAD_REQUEST
		errorData := make(map[string]interface{})
		for index, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
			errorData[strconv.Itoa(index)] = map[string]interface{}{err.Key: err.Message}
		}
		data["error"] = errorData
	}
	if AddId <= 0 {
		msg = "id 不是有效值"
		c.JSON(code, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
		c.Abort()
		return
	}
	if User.SelectAddress(AddId).UserId != userId {
		code = e.UNAUTHORIZED
		msg = "无权修改"
		c.JSON(code, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
		c.Abort()
		return
	}


	if _, ok := data["error"]; !ok {
		address := User.Address{}
		address.UserId = userId
		address.ProvinceName = c.PostForm("provinceName")
		address.CityName = c.PostForm("cityName")
		address.CountyName = c.PostForm("countyName")
		address.DetailInfo = c.PostForm("detailInfo")
		address.PostalCode = c.PostForm("postalCode")
		address.NationalCode = c.PostForm("nationalCode")
		address.UserName = c.PostForm("userName")
		address.TelNumber = c.PostForm("telNumber")
		if User.UpdateAddress(AddId, &address) {
			code = e.NO_CONTENT
		} else {
			code = e.INTERNAL_SERVER_ERROR
		}

		msg = e.GetMsg(code)
	}

	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})

}

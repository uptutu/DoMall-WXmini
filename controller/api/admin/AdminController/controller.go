package AdminController

import (
	"do-mall/models/Auth"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"do-mall/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Login(c *gin.Context) {

	data := make(map[string]interface{})
	code := e.INTERNAL_SERVER_ERROR
	var msg string

	username := c.PostForm("username")
	password := c.PostForm("password")


	valid := validation.Validation{}
	valid.Required(username, "username").Message("请输入用户名")
	valid.Required(password, "password").Message("请输入密码")

	// 处理验证错误
	if valid.HasErrors() {
		code = e.BAD_REQUEST
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		errorData := make(map[string]interface{})
		for index, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
			errorData[strconv.Itoa(index)] = map[string]interface{}{err.Key: err.Message}
		}
		// 添加错误验证错误消息到数据体
		data["error"] = errorData
	}
	if _, ok := data["error"]; !ok {

		id := Auth.CheckAdmin(username, password)
		if id > 0 {
			token, err := util.GenerateAdmin(id)
			if err == nil {
				code = e.CREATED
				data["token"] = token
			}
		} else {
			code = e.UNAUTHORIZED
			msg = "用户名或密码错误"
		}
	}

	if msg == "" {
		msg = e.GetMsg(code)
	}

	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func Show(c *gin.Context) {
	code := e.OK
	data := c.MustGet("AuthData").(*util.AdminClaims)
	c.JSON(code, gin.H{
		"code": code,
		"msg":e.GetMsg(code),
		"data": data,
	})
}
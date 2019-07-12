package UserController

import (
	"crypto/md5"
	"do-mall/models/Auth"
	"do-mall/models/User"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"do-mall/pkg/util"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"strconv"
)

//
// @Summary 注册用户
// @Produce json
// @param name query string true "mobile"
// @param state query string true "password"
// @Success 200 {string} json "{"code":201,"data":{},"msg":"Created"}"
// @Router /api/v1/user [post]
func Create(c *gin.Context) {

	data := make(map[string]interface{})
	code := e.INTERNAL_SERVER_ERROR

	mobile := c.PostForm("mobile")
	password := c.PostForm("password")

	// 数据验证
	valid := validation.Validation{}
	valid.Required(mobile, "mobile").Message("请输入手机号")
	valid.MaxSize(mobile, 11, "mobile").Message("请输入有效电话")
	valid.Phone(mobile, "mobile").Message("请输入有效电话")
	valid.Required(password, "password").Message("密码不能为空")
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
		// MD5 计算密码不可逆保存
		md5Ctx := md5.New()
		md5Ctx.Write([]byte([]byte(password)))
		password = fmt.Sprintf("%x", md5Ctx.Sum(nil))

		user := User.User{
			Mobile:   mobile,
			Password: password,
		}
		if User.CreateByPasswd(&user) {
			id := Auth.CheckAndReturnId(mobile, password)
			if id > 0 {
				token, err := util.GenerateToken(id)
				if err == nil {
					code = e.CREATED
					data["token"] = token
				}
			}

		}

	}

	c.JSON(code, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func Show(c *gin.Context) {
	code := e.OK
	data := c.MustGet("AuthData").(*util.Claims)
	c.JSON(code, gin.H{
		"code": code,
		"msg":e.GetMsg(code),
		"data": data,
	})
}

func Login(c *gin.Context)  {
	data := make(map[string]interface{})
	code := e.INTERNAL_SERVER_ERROR

	mobile := c.PostForm("mobile")
	password := c.PostForm("password")

	// 数据验证
	valid := validation.Validation{}
	valid.Required(mobile, "mobile").Message("请输入手机号")
	valid.MaxSize(mobile, 11, "mobile").Message("请输入有效电话")
	valid.Phone(mobile, "mobile").Message("请输入有效电话")
	valid.Required(password, "password").Message("密码不能为空")
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
		// MD5 计算密码不可逆保存
		md5Ctx := md5.New()
		md5Ctx.Write([]byte([]byte(password)))
		password = fmt.Sprintf("%x", md5Ctx.Sum(nil))
		// 检查登录信息
		id := Auth.CheckAndReturnId(mobile, password)
		if id > 0 {
			token, err := util.GenerateToken(id)
			if err == nil {
				code = e.CREATED
				data["token"] = token
			}
		}


	}

	c.JSON(code, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}


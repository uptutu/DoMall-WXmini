package UserController

import (
	"do-mall/models/Auth"
	"do-mall/models/User"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"do-mall/pkg/setting"
	"do-mall/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp"
	"strconv"
)

func WxLogin(c *gin.Context) {
	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})
	var msg string
	valid := validation.Validation{}

	valid.Required(c.PostForm("code"), "code").Message("code 必须")
	if valid.HasErrors() {
		code = e.BAD_REQUEST
		errorData := make(map[string]interface{})
		for index, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
			errorData[strconv.Itoa(index)] = map[string]interface{}{err.Key: err.Message}
		}
		data["error"] = errorData
	}

	if _, ok := data["error"]; ok {
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
		return
	}
	wxCode := c.PostForm("code")
	res, err := weapp.Login(setting.AppId, setting.Secret, wxCode)
	if err != nil {
		code = e.INTERNAL_SERVER_ERROR
		msg = "微信服务器登录失败"
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
		})
		c.Abort()
		return
	}
	if res.UnionID != "" {
		if user := User.QueryUserByUnionid(res.UnionID); !(user.ID > 0) {
			find := User.QueryUserByOpenid(res.OpenID)
			if find.ID > 0 {
				User.AddUnionid(find.ID, res.UnionID)
			} else {
				newUser := User.User{Openid: res.OpenID, Unionid: res.UnionID}
				if !User.Create(&newUser) {
					code = e.INTERNAL_SERVER_ERROR
					msg = "新建用户失败"
				}
			}
		}
	} else if res.OpenID != "" {
		if user := User.QueryUserByOpenid(res.OpenID); !(user.ID > 0) {
			newUser := User.User{Openid: res.OpenID, Unionid: res.UnionID}
			if !User.Create(&newUser) {
				code = e.INTERNAL_SERVER_ERROR
				msg = "新建用户失败"
			}
		}
	}
	if !User.PutSsk(res.OpenID, res.SessionKey) {
		code = e.INTERNAL_SERVER_ERROR
		msg = "存放session key失败"
	}

	id := Auth.CheckOpenid(res.OpenID)
	if id > 0 {
		token, err := util.GenerateToken(id)
		if err != nil {
			code = e.INTERNAL_SERVER_ERROR
			msg = " 查不到对应openid"
		} else {
			code = e.CREATED
			data["token"] = token
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

func WxGetUserInfo(c *gin.Context) {
	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}

	valid.Required(c.PostForm("rawData"), "rawData").Message("rawData 必须")
	valid.Required(c.PostForm("signature"), "signature").Message("signature 必须")
	valid.Required(c.PostForm("encryptedData"), "encryptedData").Message("encryptedData 必须")
	valid.Required(c.PostForm("iv"), "iv").Message("iv 必须")
	if valid.HasErrors() {
		code = e.BAD_REQUEST
		errorData := make(map[string]interface{})
		for index, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
			errorData[strconv.Itoa(index)] = map[string]interface{}{err.Key: err.Message}
		}
		data["error"] = errorData
	}

	if _, ok := data["error"]; ok {
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
		return
	}

	user := User.QueryUserByid(userId)
	rawData := c.PostForm("rawData")
	encryptedData := c.PostForm("encryptedData")
	signature := c.PostForm("signature")
	iv := c.PostForm("iv")
	ssk := user.SessionKey

	ui, err := weapp.DecryptUserInfo(rawData, encryptedData, signature, iv, ssk)
	if err != nil {
		code = e.INTERNAL_SERVER_ERROR
		msg = "解析用户信息失败"
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
		return
	}
	UserInfo := User.User{
		Avatar:   ui.Avatar,
		Nickname: ui.Nickname,
		Sex:      ui.Gender,
	}
	if !User.Update(&user, UserInfo) {
		code = e.INTERNAL_SERVER_ERROR
		msg = "更新用户信息失败"
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
		return
	}
	user = User.QueryUserByid(userId)
	data["userInfo"] = user
	code = e.OK
	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})

}

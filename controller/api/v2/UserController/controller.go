package UserController

import (
	"do-mall/models/User"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"do-mall/pkg/setting"
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

	if _, ok := data["error"]; !ok {
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
		} else {
			code = e.OK
			msg = e.GetMsg(code)
			if !User.CheckUnionid(res.UnionID) {
				User.CreateByUnionid(res.UnionID)
				if !User.PutSsk(res.UnionID, res.SessionKey) {
					code = e.INTERNAL_SERVER_ERROR
					msg = "存放session key失败"
				}
			}
		}
		c.JSON(code, gin.H{
			"openid":     res.OpenID,
			"unionId":    res.UnionID,
			"sessionKey": res.SessionKey,
		})

	}

}

func WxGetUserInfo(c *gin.Context) {
	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})
	var msg string
	valid := validation.Validation{}

	valid.Required(c.PostForm(""))
}

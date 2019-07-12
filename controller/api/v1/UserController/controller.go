package UserController

import (
	"crypto/md5"
	"do-mall/models/User"
	"do-mall/pkg/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {

	data := make(map[string]interface{})

	mobile := c.PostForm("mobile")
	password := c.PostForm("password")

	// 数据验证
	valid := validation.Validation{}
	valid.Required(mobile, "mobile").Message("请输入手机号")
	valid.MaxSize(mobile, 11, "mobile").Message("请输入有效电话")
	valid.Phone(mobile, "mobile").Message("请输入有效电话")
	valid.Required(password, "password").Message("密码不能为空")

	// MD5 计算密码不可逆保存
	md5Ctx := md5.New()
	md5Ctx.Write([]byte([]byte(password)))
	password = string(md5Ctx.Sum(nil))

	user := User.User{
		Mobile: mobile,
		Password : password,
	}
	if User.CreateByPasswd(&user) {
		code := e.SUCCESS
	} else {
		code := e.ERROR
	}

}

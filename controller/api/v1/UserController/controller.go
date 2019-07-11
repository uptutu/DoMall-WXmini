package UserController

import (
	"crypto/md5"
	"do-mall/models/User"
	"do-mall/pkg/e"
	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {

	data := make(map[string]interface{})

	mobile := c.PostForm("mobile")
	password := c.PostForm("password")

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
	}

}

package AuthController

import (
	"do-mall/models/Auth"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"do-mall/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type auth struct {
	Mobile string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	mobile := c.PostForm("mobile")
	password := c.PostForm("password")

	valid := validation.Validation{}
	a := auth{Mobile: mobile, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.BAD_REQUEST
	if ok {
		isExist := Auth.CheckAuth(mobile, password)
		if isExist {
			token, err := util.GenerateToken(mobile, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token

				code = e.OK
			}

		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : data,
	})
}

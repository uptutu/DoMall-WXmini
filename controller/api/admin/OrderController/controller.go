package OrderController

import (
	"do-mall/models/Order"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"do-mall/pkg/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Update(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}

	valid.Required(c.Param("id"), "oid").Message("oid")
	if valid.HasErrors() {
		code = e.BAD_REQUEST
		errorData := make(map[string]interface{})
		for index, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
			errorData[strconv.Itoa(index)] = map[string]interface{}{err.Key: err.Message}
		}
		data["error"] = errorData
	}
	oid := com.StrTo(c.Param("oid")).MustInt()
	order := Order.QueryOrderById(oid)
	expressTitle := c.PostForm("expressTitle")
	expressCode := c.PostForm("expressCode")
	expressNo := c.PostForm("expressNo")
	expressTime := time.Now()
	orderData := Order.Order{
		ExpressCode:  expressCode,
		ExpressNo:    expressNo,
		ExpressTitle: expressTitle,
		ExpressTime:  expressTime,
	}

	if _, ok := data["error"]; !ok {

		if order.UserId == userId {
			if Order.Update(oid, &orderData) {
				code = e.NO_CONTENT
			} else {
				code = e.INTERNAL_SERVER_ERROR
			}
		} else {
			code = e.FORBIDDEN
		}

	}

	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

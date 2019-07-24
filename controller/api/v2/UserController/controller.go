package UserController

import (
	"do-mall/models/Auth"
	"do-mall/models/User"
	"do-mall/models/Wallet"
	"do-mall/models/WxOrder"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"do-mall/pkg/setting"
	"do-mall/pkg/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp"
	"github.com/medivhzhan/weapp/payment"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"strconv"
)

func WxLogin(c *gin.Context) {
	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})
	var msg string
	valid := validation.Validation{}

	valid.Required(c.Param("code"), "code").Message("code 必须")
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

func WxGetPhone(c *gin.Context) {
	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}

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
	encryptedData := c.PostForm("encryptedData")
	iv := c.PostForm("iv")
	ssk := user.SessionKey

	phone, err := weapp.DecryptPhoneNumber(ssk, encryptedData, iv)
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
	if !User.UpdateColumn(&user, "mobile", phone.PurePhoneNumber) {
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

func WxTopup(c *gin.Context) {
	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}

	valid.Required(c.PostForm("amount"), "amount").Message("amount 必须")
	valid.Numeric(c.PostForm("amount"), "amount").Message("amount 必须是有效数值,单位为分")
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

	amount := com.StrTo(c.PostForm("amount")).MustInt()
	sumPay := float64(amount / 100)
	tradeNo := string(com.RandomCreateBytes(32))
	receiveUrl := "https://xxx.xxx.xxx"
	Body := "充值 " + com.ToStr(sumPay) + "RMB"
	user := User.QueryUserByid(userId)


	form := payment.Order{
		AppID:      setting.AppId,
		MchID:      setting.MchID,
		Body:       Body,
		NotifyURL:  receiveUrl,
		OpenID:     user.Openid,
		OutTradeNo: tradeNo,
		TotalFee:   amount,
		Detail:     Body,
		Attach:     Body,
	}

	res, err := form.Unify(setting.PayKey)
	if err != nil {
		logging.Info(err)
		code = e.INTERNAL_SERVER_ERROR
		msg = fmt.Sprint("出错问题: %v", err)
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
		return
	}

	params, err := payment.GetParams(res.AppID, setting.PayKey, res.NonceStr, res.PrePayID)
	if err != nil {
		logging.Info(err)
		code = e.INTERNAL_SERVER_ERROR
		msg = fmt.Sprint("出错问题: %v", err)
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
		return
	}
	order := WxOrder.WxOrder{
		UserId:userId,
		OutTradeNo:tradeNo,
		NonceStr:res.NonceStr,
		Sign:res.Sign,
		Body:Body,
		Detail:form.Detail,
		Attach:form.Attach,
		SumPay:sumPay,
		TotalFee:form.TotalFee,
	}
	if !WxOrder.Create(&order) {
		logging.Info("未能插入微信订单表")
		code = e.INTERNAL_SERVER_ERROR
		msg = fmt.Sprint("未能插入微信订单表")
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
		return
	}
	code = e.OK
	msg = e.GetMsg(code)
	data["data"] = params
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func WxTopupCallback(c *gin.Context) {
	err := payment.HandlePaidNotify(c.Writer, c.Request,  func(ntf payment.PaidNotify) (bool, string) {
		tradeNo := ntf.OutTradeNo
		order := WxOrder.QueryByTradeNo(tradeNo)
		if Wallet.TopUpBalance(order.SumPay, order.UserId) {
			order.TransactionId = ntf.TransactionID
			if WxOrder.Update(&order) {
				return true, ""
			}
		}
		return false, "服务器内部更新错误"
	})
	if err != nil {
		logging.Info(err)
	}
}
package OrderController

import (
	"do-mall/models/Order"
	"do-mall/models/Product"
	"do-mall/models/User"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"do-mall/pkg/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Create(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}

	valid.Required(c.PostForm("cartIds"), "cartIds").Message("cartId 必须")
	cartIds := c.PostFormArray("cartIds")
	for _, cartId := range cartIds {
		valid.Numeric(cartId, "cartId").Message("cartId 必须是有效数值")
	}

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
		var address User.Address
		var total float64
		if addressId := com.StrTo(c.PostForm("addressId")).MustInt(); addressId != 0 {
			address = User.SelectAddress(addressId)
		} else {
			address = User.SelectDefaultAddress(userId)
		}
		carts := User.QueryCarts(cartIds)
		// *** Can be optimized
		for _, c := range carts {
			total += float64(c.Number) * Product.Show(c.PId).Price
		}

		order := Order.Order{
			UserId:       userId,
			ProvinceName: address.ProvinceName,
			CityName:     address.CityName,
			CountyName:   address.CountyName,
			DetailInfo:   address.DetailInfo,
			PostalCode:   address.PostalCode,
			UserName:     address.UserName,
			TelNumber:    address.TelNumber,
			Total:        total,
			SumPay:       total,
		}

		if Order.Created(&order) {
			if User.OrderCreated(cartIds, order.ID) {
				code = e.CREATED
			}
		} else {
			code = e.INTERNAL_SERVER_ERROR
		}
	}

	msg = e.GetMsg(code)
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
	valid := validation.Validation{}

	valid.Required(c.Param("oid"), "oid").Message("oid 必须")
	valid.Numeric(c.Param("oid"), "oid").Message("oid 必须是有效数字")

	oid := com.StrTo(c.Param("oid")).MustInt()
	order := Order.QueryOrderById(oid)
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
		if order.UserId == userId {
			if Order.Destroy(oid) {
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

func List(c *gin.Context) {
	code := e.OK
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID

	data["orders"] = Order.QueryOrderByUserId(userId)
	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func Pay(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID

	valid := validation.Validation{}

	valid.Required(c.Param("oid"), "oid").Message("oid 必须")
	valid.Numeric(c.Param("oid"), "oid").Message("oid 必须是有效数字")

	oid := com.StrTo(c.Param("oid")).MustInt()
	order := Order.QueryOrderById(oid)
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
		if order.UserId == userId {
			if Order.Settlement(userId, oid) {
				code = e.NO_CONTENT
				msg = e.GetMsg(code)
			} else {
				code = e.FORBIDDEN
				msg = "余额不足"
			}
		}
		code = e.FORBIDDEN
		msg = e.GetMsg(code)
	}
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func ViewOrderDetails(c *gin.Context)  {
	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}
	valid.Required(c.Param("id"), "id").Message("id(Order Id) 必须")
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
	oid := com.StrTo(c.Param("id")).MustInt()
	if !Order.IsOwner(userId, oid) {
		code = e.FORBIDDEN
		msg = e.GetMsg(code)
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		c.Abort()
		return
	}
	carts := User.QueryCartByOid(oid)
	for i, cart := range carts {
		data[string(i)] = Product.Show(cart.PId)
	}
	code = e.OK
	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
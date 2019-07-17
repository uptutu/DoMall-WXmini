package ProductController

import (
	"do-mall/models/Product"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
	"strconv"
)

func Create(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	valid := validation.Validation{}

	product := Product.Product{}
	product.Title = c.PostForm("title")
	product.Cover = c.PostForm("cover")
	product.Carousel = c.PostForm("carousel")
	product.Brand = c.PostForm("brand")
	product.Series = c.PostForm("series")
	product.Price = com.StrTo(c.PostForm("price")).MustFloat64()
	product.SellingPrice = com.StrTo(c.PostForm("selling_price")).MustFloat64()
	product.Cost = com.StrTo(c.PostForm("cost")).MustFloat64()
	product.Tags = c.PostForm("tags")
	product.Inventory = com.StrTo(c.PostForm("sales")).MustInt()
	product.Status = com.StrTo(c.PostForm("inventory")).MustInt()
	product.OnSale = com.StrTo(c.PostForm("status")).MustInt()

	log.Println(product)

	valid.Required(product.Title, "Title").Message("请输入商品名")
	valid.MaxSize(product.Title, 80, "Title").Message("请输入商品名")

	valid.Required(product.Cover, "Cover").Message("请出入封面图")
	valid.Match(product.Cover, regexp.MustCompile(`(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`), "Cover").Message("请出入封面图")

	valid.Required(product.Carousel, "Carousel").Message("请输入图集")

	valid.Required(product.Brand, "Brand").Message("请输入商品品牌")

	valid.Required(product.Series, "Series").Message("请输入商品系列名")

	valid.Required(product.Price, "Price").Message("请输入商品零售价")
	valid.Required(product.SellingPrice, "SellingPrice").Message("请输入商品售价")
	valid.Required(product.Cost, "Cost").Message("请输入商品成本价")

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
		if Product.Create(&product) {
			code = e.CREATED
		}
	} else {
		code = e.INTERNAL_SERVER_ERROR
	}

	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func Update(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}

	product := Product.Product{}
	product.ID = id
	product.Title = c.PostForm("title")
	product.Cover = c.PostForm("cover")
	product.Carousel = c.PostForm("carousel")
	product.Brand = c.PostForm("brand")
	product.Series = c.PostForm("series")
	product.Price = com.StrTo(c.PostForm("price")).MustFloat64()
	product.SellingPrice = com.StrTo(c.PostForm("selling_price")).MustFloat64()
	product.Cost = com.StrTo(c.PostForm("cost")).MustFloat64()
	product.Tags = c.PostForm("tags")

	product.Sales = com.StrTo(c.PostForm("sales")).MustInt()
	product.Status = com.StrTo(c.PostForm("status")).MustInt()
	product.OnSale = com.StrTo(c.PostForm("on_sale")).MustInt()

	valid.Required(product.ID, "Id").Message("请输入Id")
	valid.Min(product.ID, 0, "Id").Message("请输入有效id")

	valid.Required(product.Title, "Title").Message("请输入商品名")
	valid.MaxSize(product.Title, 80, "Title").Message("请输入商品名")

	valid.Required(product.Cover, "Cover").Message("请出入封面图")
	valid.Match(product.Cover, regexp.MustCompile(`(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`), "Cover").Message("请出入封面图")

	valid.Required(product.Carousel, "Carousel").Message("请输入图集")

	valid.Required(product.Brand, "Brand").Message("请输入商品品牌")

	valid.Required(product.Series, "Series").Message("请输入商品系列名")

	valid.Required(product.Price, "Price").Message("请输入商品零售价")
	valid.Required(product.SellingPrice, "SellingPrice").Message("请输入商品售价")
	valid.Required(product.Cost, "Cost").Message("请输入商品成本价")

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
		if Product.Update(&product) {
			code = e.NO_CONTENT
		}
	} else {
		code = e.INTERNAL_SERVER_ERROR
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

	id := com.StrTo(c.Param("id")).MustInt()
	if Product.Destroy(id) {
		code = e.NO_CONTENT
	} else {
		code = e.INTERNAL_SERVER_ERROR
	}

	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func UpdateInventory(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	valid := validation.Validation{}

	pid := com.StrTo(c.Param("pid")).MustInt()

	inv := Product.Inventory{}
	inv.F = com.StrTo(c.PostForm("f")).MustInt()
	valid.Numeric(inv.F, "f").Message("请输有效整数数字")
	inv.Xxs = com.StrTo(c.PostForm("xxs")).MustInt()
	valid.Numeric(inv.Xxs, "xss").Message("请输有效整数数字")
	inv.Xs = com.StrTo(c.PostForm("xs")).MustInt()
	valid.Numeric(inv.Xs, "xs").Message("请输有效整数数字")
	inv.S = com.StrTo(c.PostForm("s")).MustInt()
	valid.Numeric(inv.S, "s").Message("请输有效整数数字")
	inv.M = com.StrTo(c.PostForm("m")).MustInt()
	valid.Numeric(inv.M, "m").Message("请输有效整数数字")
	inv.L = com.StrTo(c.PostForm("l")).MustInt()
	valid.Numeric(inv.L, "l").Message("请输有效整数数字")
	inv.Xl = com.StrTo(c.PostForm("xl")).MustInt()
	valid.Numeric(inv.Xl, "xl").Message("请输有效整数数字")
	inv.Xxl = com.StrTo(c.PostForm("xxl")).MustInt()
	valid.Numeric(inv.Xxl, "xxl").Message("请输有效整数数字")
	inv.S35 = com.StrTo(c.PostForm("s35")).MustInt()
	valid.Numeric(inv.S35, "s35").Message("请输有效整数数字")
	inv.S36 = com.StrTo(c.PostForm("s36")).MustInt()
	valid.Numeric(inv.S36, "s36").Message("请输有效整数数字")
	inv.S37 = com.StrTo(c.PostForm("s37")).MustInt()
	valid.Numeric(inv.S37, "s37").Message("请输有效整数数字")
	inv.S38 = com.StrTo(c.PostForm("s38")).MustInt()
	valid.Numeric(inv.S38, "s38").Message("请输有效整数数字")
	inv.S39 = com.StrTo(c.PostForm("s39")).MustInt()
	valid.Numeric(inv.S39, "s39").Message("请输有效整数数字")
	inv.S40 = com.StrTo(c.PostForm("s40")).MustInt()
	valid.Numeric(inv.S40, "s40").Message("请输有效整数数字")
	inv.S405 = com.StrTo(c.PostForm("s405")).MustInt()
	valid.Numeric(inv.S405, "s405").Message("请输有效整数数字")
	inv.S41 = com.StrTo(c.PostForm("s41")).MustInt()
	valid.Numeric(inv.S41, "s41").Message("请输有效整数数字")
	inv.S415 = com.StrTo(c.PostForm("s415")).MustInt()
	valid.Numeric(inv.S415, "s415").Message("请输有效整数数字")
	inv.S42 = com.StrTo(c.PostForm("s42")).MustInt()
	valid.Numeric(inv.S42, "s42").Message("请输有效整数数字")
	inv.S425 = com.StrTo(c.PostForm("s425")).MustInt()
	valid.Numeric(inv.S425, "s425").Message("请输有效整数数字")
	inv.S43 = com.StrTo(c.PostForm("s43")).MustInt()
	valid.Numeric(inv.S43, "s43").Message("请输有效整数数字")
	inv.S435 = com.StrTo(c.PostForm("s435")).MustInt()
	valid.Numeric(inv.S435, "s435").Message("请输有效整数数字")
	inv.S44 = com.StrTo(c.PostForm("s44")).MustInt()
	valid.Numeric(inv.S44, "s44").Message("请输有效整数数字")
	inv.S445 = com.StrTo(c.PostForm("s445")).MustInt()
	valid.Numeric(inv.S445, "s445").Message("请输有效整数数字")
	inv.S45 = com.StrTo(c.PostForm("s45")).MustInt()
	valid.Numeric(inv.S45, "s45").Message("请输有效整数数字")
	inv.S46 = com.StrTo(c.PostForm("s46")).MustInt()
	valid.Numeric(inv.S46, "s46").Message("请输有效整数数字")
	inv.S47 = com.StrTo(c.PostForm("s47")).MustInt()
	valid.Numeric(inv.S47, "s47").Message("请输有效整数数字")

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
		if Product.UpdateInventory(pid, &inv) {
			code = e.NO_CONTENT
		}
	} else {
		code = e.INTERNAL_SERVER_ERROR
	}

	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})

}

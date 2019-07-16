package ProductController

import (
	"do-mall/models/Product"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
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

	valid.Required(product.Title, "Title").Message("请输入商品名")
	valid.MaxSize(product.Title, 80,"Title").Message("请输入商品名")

	valid.Required(product.Cover, "Cover").Message("请出入封面图")
	valid.Match(product.Cover, regexp.MustCompile(`(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`), "Cover").Message("请出入封面图")

	valid.Required(product.Carousel, "Carousel").Message("请输入图集")

	valid.Required(product.Brand, "Brand").Message("请输入商品品牌")

	valid.Required(product.Series, "Series").Message("请输入商品系列名")

	valid.Required(product.Price, "Price").Message("请输入商品零售价")
	valid.Min(product.Price, 0, "Price").Message("商品零售价不得小于0元")

	valid.Required(product.SellingPrice, "SellingPrice").Message("请输入商品售价")
	valid.Min(product.SellingPrice, 0, "SellingPrice").Message("商品售价不得小于0元")

	valid.Required(product.Cost, "Cost").Message("请输入商品成本价")
	valid.Min(product.Cost, 0, "Cost").Message("商品成本价不得小于0元")

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
	product.Inventory = com.StrTo(c.PostForm("sales")).MustInt()
	product.Status = com.StrTo(c.PostForm("inventory")).MustInt()
	product.OnSale = com.StrTo(c.PostForm("status")).MustInt()

	valid.Required(product.ID, "Id").Message("请输入Id")
	valid.Min(product.ID, 0, "Id").Message("请输入有效id")

	valid.Required(product.Title, "Title").Message("请输入商品名")
	valid.MaxSize(product.Title, 80,"Title").Message("请输入商品名")

	valid.Required(product.Cover, "Cover").Message("请出入封面图")
	valid.Match(product.Cover, regexp.MustCompile(`(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`), "Cover").Message("请出入封面图")

	valid.Required(product.Carousel, "Carousel").Message("请输入图集")

	valid.Required(product.Brand, "Brand").Message("请输入商品品牌")

	valid.Required(product.Series, "Series").Message("请输入商品系列名")

	valid.Required(product.Price, "Price").Message("请输入商品零售价")
	valid.Min(product.Price, 0, "Price").Message("商品零售价不得小于0元")

	valid.Required(product.SellingPrice, "SellingPrice").Message("请输入商品售价")
	valid.Min(product.SellingPrice, 0, "SellingPrice").Message("商品售价不得小于0元")

	valid.Required(product.Cost, "Cost").Message("请输入商品成本价")
	valid.Min(product.Cost, 0, "Cost").Message("商品成本价不得小于0元")

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

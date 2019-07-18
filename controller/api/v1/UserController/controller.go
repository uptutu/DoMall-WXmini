package UserController

import (
	"crypto/md5"
	"do-mall/models/Auth"
	"do-mall/models/User"
	"do-mall/pkg/e"
	"do-mall/pkg/logging"
	"do-mall/pkg/setting"
	"do-mall/pkg/util"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
)

//
// @Summary 注册用户
// @Produce json
// @param name query string true "mobile"
// @param state query string true "password"
// @Success 200 {string} json "{"code":201,"data":{},"msg":"Created"}"
// @Router /api/v1/user [post]
func Create(c *gin.Context) {

	data := make(map[string]interface{})
	code := e.INTERNAL_SERVER_ERROR
	var msg string

	mobile := c.PostForm("mobile")
	password := c.PostForm("password")

	valid := validation.Validation{}
	valid.Required(mobile, "mobile").Message("请输入手机号")
	valid.MaxSize(mobile, 11, "mobile").Message("请输入有效电话")
	valid.Phone(mobile, "mobile").Message("请输入有效电话")
	valid.Required(password, "password").Message("密码不能为空")

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
		md5Ctx := md5.New()
		md5Ctx.Write([]byte([]byte(password)))
		password = fmt.Sprintf("%x", md5Ctx.Sum(nil))

		user := User.User{
			Mobile:   mobile,
			Password: password,
		}
		if User.CreateByPasswd(&user) {
			id := Auth.CheckAndReturnId(mobile, password)
			if id > 0 {
				token, err := util.GenerateToken(id)
				if err == nil {
					code = e.CREATED
					data["token"] = token
				}
			} else {
				code = e.UNAUTHORIZED
				msg = "用户名或密码错误"
			}

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

func Show(c *gin.Context) {
	code := e.OK
	data := make(map[string]interface{})
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	data["user"] = User.GetInfo(userId)

	c.JSON(code, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func Login(c *gin.Context) {
	data := make(map[string]interface{})
	code := e.INTERNAL_SERVER_ERROR
	var msg string

	mobile := c.PostForm("mobile")
	password := c.PostForm("password")

	// 数据验证
	valid := validation.Validation{}
	valid.Required(mobile, "mobile").Message("请输入手机号")
	valid.MaxSize(mobile, 11, "mobile").Message("请输入有效电话")
	valid.Phone(mobile, "mobile").Message("请输入有效电话")
	valid.Required(password, "password").Message("密码不能为空")
	// 处理验证错误
	if valid.HasErrors() {
		code = e.BAD_REQUEST
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		errorData := make(map[string]interface{})
		for index, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
			errorData[strconv.Itoa(index)] = map[string]interface{}{err.Key: err.Message}
		}
		// 添加错误验证错误消息到数据体
		data["error"] = errorData
	}
	if _, ok := data["error"]; !ok {
		// MD5 计算密码不可逆保存
		md5Ctx := md5.New()
		md5Ctx.Write([]byte([]byte(password)))
		password = fmt.Sprintf("%x", md5Ctx.Sum(nil))
		// 检查登录信息
		id := Auth.CheckAndReturnId(mobile, password)
		if id > 0 {
			token, err := util.GenerateToken(id)
			if err == nil {
				code = e.CREATED
				data["token"] = token
			}
		} else {
			code = e.UNAUTHORIZED
			msg = "用户名或密码错误"
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

func Update(c *gin.Context) {
	// 初始返回数据
	var msg string
	code := e.INTERNAL_SERVER_ERROR
	data := make(map[string]interface{})

	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	nickname := c.PostForm("nickname")
	avatar := c.PostForm("avatar")
	sex := c.PostForm("sex")
	introduction := c.PostForm("introduction")

	valid := validation.Validation{}
	editedData := make(map[string]interface{})

	if nickname != "" {
		err := valid.MaxSize(nickname, 10, "nickname").Message("限定10个字符").Error
		if err == nil {
			editedData["nickname"] = nickname
		}
	}

	if avatar != "" {
		reg := regexp.MustCompile(`(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`)
		err := valid.Match(avatar, reg, "avatar").Message("请上传正确图片").Error
		if err == nil {
			editedData["avatar"] = avatar
		}
	}

	if sex != "" {
		err := valid.Numeric(sex, "sex").Message("请传入有效数据").Error
		sexInt, _ := com.StrTo(sex).Int()
		errMin := valid.Min(sexInt, 0, "sex").Message("性别传值不正确").Error
		errMax := valid.Max(sexInt, 2, "sex").Message("性别传值不正确").Error
		if errMin == nil && errMax == nil && err != nil {
			editedData["sex"] = sexInt
		}
	}

	if introduction != "" {
		err := valid.MaxSize(introduction, 80, "introduction").Message("限定80个字符").Error
		if err == nil {
			editedData["introduction"] = introduction
		}
	}

	if valid.HasErrors() {
		code = e.BAD_REQUEST
		errorData := make(map[string]interface{})
		for index, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
			errorData[strconv.Itoa(index)] = map[string]interface{}{err.Key: err.Message}
		}
		data["error"] = errorData
		c.JSON(code, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
		c.Abort()
		return
	}
	user := User.GetInfo(userId)

	// 更新数据
	if User.Update(&user, editedData) {
		code = e.OK
		msg = "更新成功"
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

func FavoritesList(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID

	data["lists"], data["total"] = User.ShowFavorites(util.GetPage(c), setting.PageSize, userId)
	if _, ok := data["lists"]; ok {
		code = e.OK
	}
	msg = e.GetMsg(code)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func FavoritesCreate(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}

	valid.Required(c.PostForm("pId"), "pId").Message("pId 必须")
	valid.Numeric(c.PostForm("pId"), "pId").Message("pId 必须是有效数值")

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
		pId := com.StrTo(c.PostForm("pId")).MustInt()
		if User.AddFavorite(userId, pId) {
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

func FavoritesDestroy(c *gin.Context) {
	code := e.BAD_REQUEST
	data := make(map[string]interface{})
	var msg string
	userId := c.MustGet("AuthData").(*util.Claims).User.ID
	valid := validation.Validation{}

	valid.Required(c.Param("id"), "pId").Message("pId 参数值必须")
	valid.Numeric(c.Param("id"), "pId").Message("pId 必须是有效数值")

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
		pId := com.StrTo(c.Param("id")).MustInt()
		if User.DestroyFavorite(userId, pId) {
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

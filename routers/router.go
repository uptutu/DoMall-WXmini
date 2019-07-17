package routers

import (
	"do-mall/controller/api/admin/AdminController"
	ProductController2 "do-mall/controller/api/admin/ProductController"
	"do-mall/controller/api/v1/ProductController"
	"do-mall/controller/api/v1/UserController"
	_ "do-mall/docs" // docs is generated by Swag CLI, you have to import it.
	jwt "do-mall/middleware"
	"do-mall/pkg/e"
	"do-mall/pkg/setting"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())


	gin.SetMode(setting.RunMode)

	// use ginSwagger middleware to
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册路由
	apiv1 := r.Group("/api/v1")
	{
		// User
		userRoute := apiv1.Group("/user")
		{
			userRoute.POST("/register", UserController.Create)
			userRoute.POST("/login", UserController.Login)
			userRoute.GET("/index", jwt.JWT(), UserController.Show)
			userRoute.PUT("/edit/:id", jwt.JWT(), UserController.Update)

			// Favorites
			userRoute.GET("/favorite", jwt.JWT(), UserController.FavoritesList)
			userRoute.POST("/favorite", jwt.JWT(), UserController.FavoritesCreate)
			userRoute.DELETE("/favorite", jwt.JWT(), UserController.FavoritesDestroy)
		}

		// Products
		productRoute := apiv1.Group("/product")
		{
			productRoute.GET("/index", ProductController.Index)
			productRoute.GET("/search", ProductController.Search)
			productRoute.GET("/info/:id", ProductController.Show)

		}



	}

	// 注册后台路由
	r.POST("/adminUserLogin", AdminController.Login)
	admin := r.Group("/admin")
	admin.Use(jwt.Admin())
	{
		// Admin Info
		admin.GET("/info", AdminController.Show)

		// Products
		adminProductRoute := admin.Group("/product")
		{
			adminProductRoute.POST("/", ProductController2.Create)
			adminProductRoute.PUT("/:id", ProductController2.Update)
			adminProductRoute.DELETE("/:id", ProductController2.Destroy)
		}

	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":404,
			"msg":e.GetMsg(404),
			"data":nil,
		})
	})

	return r
}

package router

import (
	"go_web_bbs/controller"
	"go_web_bbs/logger"
	"go_web_bbs/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// use ginSwagger middleware to
	url := ginSwagger.URL("http://localhost:8081/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	//r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	v1 := r.Group("/api/v1")
	{
		// 注册
		v1.POST("/signup", controller.SignUpHandler)
		// 登录
		v1.POST("/login", controller.LoginHandler)
		// 应用JWT认证中间件
		v1.Use(middlewares.JWTAuthMiddleware())

		community := v1.Group("/community")
		{
			community.GET("/", controller.CommunityHandler)
			community.GET("/:id", controller.CommunityDetailHandler)
		}

		post := v1.Group("/post")
		{
			post.POST("/", controller.CreatePostHandler)
			//post.GET("/:id", controller.PostDetailHandler)
			//post.GET("/", controller.PostListHandler)
		}

	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}

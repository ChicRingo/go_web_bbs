package router

import (
	"fmt"
	"go_web_bbs/controller"
	_ "go_web_bbs/docs"
	"go_web_bbs/logger"
	"go_web_bbs/middlewares"
	"go_web_bbs/settings"
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

	// 使用 ginSwagger 中间件
	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", settings.Conf.Port)) // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 第一版API接口
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
			post.GET("/:id", controller.GetPostDetailHandler)
		}
		v1.GET("/posts", controller.GetPostListHandler)
		// 根据时间或者分数获取帖子列表
		v1.GET("/posts2", controller.GetPostListHandler2)
		// 根据社区获取帖子列表
		v1.GET("/posts3", controller.GetCommunityPostListHandler)

		vote := v1.Group("/vote")
		{
			vote.POST("/", controller.PostVoteHandler)
		}

	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}

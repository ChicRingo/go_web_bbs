package controller

import (
	"go_web_bbs/logic"
	"go_web_bbs/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler godoc
// @Summary 创建帖子
// @Description 根据get请求参数创建帖子
// @Tags post
// @version 1.0
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param post body models.Post true "提交文章对象"
// @Success 1000 "success" {object} controller.ResponseData
// @Failure 1001 "请求参数错误" {object} controller.ResponseData
// @Failure 1005 "服务繁忙" {object} controller.ResponseData
// @Failure 1006 "需要登录" {object} controller.ResponseData
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}
	// 参数校验

	// 获取作者ID，当前请求的UserID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	post.AuthorId = userID

	err = logic.CreatePost(&post)
	if err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// PostListHandler godoc
// @Summary 帖子列表
// @Description 获取全部帖子列表
// @Tags post
// @version 1.0
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param order query int true "每页数量"
// @Param page query int true "分页页码"
// @Success 1000 "success" {object} controller.ResponseData
// @Failure 1001 "请求参数错误" {object} controller.ResponseData
// @Failure 1005 "服务繁忙" {object} controller.ResponseData
// @Failure 1006 "需要登录" {object} controller.ResponseData
// @Router /post [get]
//func PostListHandler(c *gin.Context) {
//	order, _ := c.GetQuery("order")
//	pageStr, ok := c.GetQuery("page")
//	if !ok {
//		pageStr = "1"
//	}
//	pageNum, err := strconv.ParseInt(pageStr, 10, 64)
//	if err != nil {
//		pageNum = 1
//	}
//	posts := redis.GetPost(order, pageNum)
//	fmt.Println(len(posts))
//	ResponseSuccess(c, posts)
//}

// PostDetailHandler 帖子详情
//func PostDetailHandler(c *gin.Context) {
//	postId := c.Param("id")
//
//	post, err := logic.GetPost(postId)
//	if err != nil {
//		zap.L().Error("logic.GetPost(postID) failed", zap.String("postId", postId), zap.Error(err))
//	}
//
//	ResponseSuccess(c, post)
//}

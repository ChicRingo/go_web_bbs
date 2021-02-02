package controller

import (
	"go_web_bbs/logic"
	"go_web_bbs/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler godoc
// @Summary 创建帖子
// @Description 根据get请求参数创建帖子
// @Tags post
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param post body models.Post true "提交文章对象"
// @Success 1000 {object} controller.ResponseData
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	// 1.获取参数及参数校验
	post := new(models.Post)
	if err := c.ShouldBindJSON(post); err != nil {
		zap.L().Debug("c.ShouldBindJSON(post) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}

	// 2.获取作者ID (从 c 取到当前请求的UserID)
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	post.AuthorID = userID

	// 3.创建帖子
	if err = logic.CreatePost(post); err != nil {
		zap.L().Error("logic.CreatePost(post) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 4.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler godoc
// @Summary 帖子详情
// @Description 根据id获取帖子详情
// @Tags post
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path int true "post_id"
// @Success 1000 {object} controller.ResponseData
// @Failure 1001 {object} controller.ResponseData
// @Failure 1005 {object} controller.ResponseData
// @Failure 1006 {object} controller.ResponseData
// @Router /post/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	// 1.获取参数及参数校验
	postIdStr := c.Param("id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2.根据id取出帖子数据（查询数据库）
	post, err := logic.GetPostById(postId)
	if err != nil {
		zap.L().Error("logic.GetPostById(postID) failed", zap.Int64("postId", postId), zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	// 3.成功返回响应
	ResponseSuccess(c, post)
}

// GetPostListHandler godoc
// @Summary 帖子列表
// @Description 获取全部帖子列表
// @Tags post
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param page query int true "分页页码"
// @Param size query int true "每页数量"
// @Success 1000 {object} controller.ResponseData
// @Failure 1005 {object} controller.ResponseData
// @Router /posts1 [get]
func GetPostListHandler(c *gin.Context) {
	// 1.获取参数及参数校验
	page, size := getPageInfo(c)
	// 2.根据页码和个数获取分页数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 godoc
// @Summary 帖子列表
// @Description 根据前端参数动态获取全部帖子列表，按 时间 或 分数 排序
// @Tags post
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string false "Bearer 用户令牌"
// @Param page query int true "分页页码"
// @Param size query int true "每页数量"
// @Param order query string true "排序规则"
// @Success 1000 {object} controller.ResponseData
// @Failure 1005 {object} controller.ResponseData
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// 1.获取参数及参数校验

	// 初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 3.根据id、页码和个数获取分页数据
	data, err := logic.GetPostListByOrder(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, data)
}

// GetCommunityPostListHandler godoc
// @Summary 帖子列表
// @Description 根据社区获取帖子列表
// @Tags post
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param page query int true "分页页码"
// @Param size query int true "每页数量"
// @Param order query string true "排序规则"
// @Param community_id query string true "社区id"
// @Success 1000 {object} controller.ResponseData
// @Failure 1005 {object} controller.ResponseData
// @Router /posts2 [get]
func GetCommunityPostListHandler(c *gin.Context) {
	// 初始化结构体时指定初始参数
	p := &models.ParamCommunityList{
		ParamPostList: models.ParamPostList{
			Page:  1,
			Size:  10,
			Order: models.OrderTime,
		},
	}

	// 绑定 query 参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 根据id、页码和个数获取分页数据
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		zap.L().Error("logic.GetCommunityPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, data)
}

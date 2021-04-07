package controller

import (
	"errors"
	"go_web_bbs/dao/mysql"
	"go_web_bbs/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区

// CommunityHandler godoc
// @Summary 社区列表
// @Description 查询到所有的社区（community_id, community_name）以列表形式返回
// @Tags community
// @Security ApiKeyAuth
// @Success 200 {object} _responseCommunityList
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	// 1.获取全部社区列表
	communityList, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("mysql.GetCommunityList() failed", zap.Error(err))

		//不轻易把服务器错误暴露给外界
		ResponseError(c, CodeServerBusy)
		return
	}

	// 2.返回社区列表
	ResponseSuccess(c, communityList)
}

// CommunityDetailHandler godoc
// @Summary 社区分类详情
// @Description 根据id获取社区详情
// @Tags community
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path int true "社区ID"
// @Success 200 {object} _responseCommunityDetail
// @Router /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	// 1.URL路径参数获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// 参数错误
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2.根据id获取社区详情
	communityDetail, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID() failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorInvalidID) {
			ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
			return
		}

		ResponseError(c, CodeServerBusy) //不轻易把服务器错误暴露给外界
		return
	}

	// 3.返回社区详情
	ResponseSuccess(c, communityDetail)
}

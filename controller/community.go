package controller

import (
	"fmt"
	"go_web_bbs/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区

// CommunityHandler 社区列表
// 查询到所有的社区（community_id, community_name）以列表形式返回
func CommunityHandler(c *gin.Context) {
	communityList, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("mysql.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务器错误暴露给外界
		return
	}
	ResponseSuccess(c, communityList)
}

// GetCommunityDetail 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	// 1.获取社区id
	idStr := c.Param("id") // 获取URL路径参数
	fmt.Println(idStr)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam) // 参数错误
		return
	}
	//2. 根据id获取社区详情
	communityDetail, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务器错误暴露给外界
		return
	}
	fmt.Println(communityDetail)
	ResponseSuccess(c, communityDetail)
}

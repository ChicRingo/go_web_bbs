package controller

import (
	"errors"
	"fmt"
	"go_web_bbs/dao/redis"
	"go_web_bbs/logic"
	"go_web_bbs/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票

//type VoteData struct {
//	// UserID 从请求中获取当前的用户
//	PostID    int64 `json:"post_id,string"`   // 贴子id
//	Direction int   `json:"direction,string"` // 赞成票(1)还是反对票(-1)
//}
func PostVoteHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		valErrs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errMsg := removeTopStruct(valErrs.Translate(trans)) // 翻译并去除掉错误提示中的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParam, errMsg)
		return
	}
	// 获取当前请求的用户的id
	userID, err := getCurrentUserID(c)
	fmt.Println(userID)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 具体投票的业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		if errors.Is(err, redis.ErrVoteRepeated) {
			ResponseSuccess(c, err.Error())
			return
		}
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

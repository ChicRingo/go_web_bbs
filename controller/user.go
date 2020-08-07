package controller

import (
	"errors"
	"go_web_bbs/dao/mysql"
	"go_web_bbs/logic"
	"go_web_bbs/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// @Param article body model.Article true "文章"
// @Success 200 object model.Result 成功后返回值
// @Failure 409 object model.Result 添加失败
// @Router /article [post]

// @Summary 处理注册请求的函数
// @Id 1
// @Tags 注册用户
// @version 1.0
// @Accept application/x-json-stream
// @Param paramSingUp body models.ParamSingUp true "用户注册请求"
// @Success 1000 object 成功后返回值
// @Failure 409 object 添加失败
// @Router /v1/api/signup [post]
func SignUpHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	p := new(models.ParamSingUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator。ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}

	// 2.业务处理
	if err := logic.SingUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, nil)
}

// SignUpHandler 处理注册请求的函数
func LoginHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator。ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}
	// 2.业务处理
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, token)
}

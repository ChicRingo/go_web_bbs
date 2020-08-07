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

// SignUpHandler godoc
// @Summary 新建用户
// @Description 根据传递进来的用户名和密码进行校验，通过后创建新用户
// @Tags user
// @version 1.0
// @Accept json
// @Produce json
// @Param paramSingUp body models.ParamSingUp true "用户注册请求"
// @Success 1000 "success" {object} controller.ResponseData
// @Failure 1001 "请求参数错误" {object} controller.ResponseData
// @Failure 1005 "服务繁忙" {object} controller.ResponseData
// @Failure 1002 "用户名已存在" {object} controller.ResponseData
// @Router /signup [post]
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
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
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

// LoginHandler godoc
// @Summary 用户登录
// @Description 根据传递进来的用户名和密码进行校验，通过后创建新用户
// @Tags user
// @version 1.0
// @Accept json
// @Produce json
// @Param paramLogin body models.ParamLogin true "用户登陆请求"
// @Success 1000 "success" {object} controller.ResponseData
// @Failure 1001 "请求参数错误" {object} controller.ResponseData
// @Failure 1005 "服务繁忙" {object} controller.ResponseData
// @Failure 1003 "用户名不存在" {object} controller.ResponseData
// @Router /login [post]
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
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
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

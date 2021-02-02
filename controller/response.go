package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": 10000, // 程序中的错误码
	"msg": xx,     // 提示信息
	"data": {},    // 数据
}
*/

type ResponseData struct {
	Code ResCode     `json:"code" example:"1000"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 请求响应成功，返回成功信息
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

// 请求响应错误，返回错误信息
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// 请求响应错误，返回错误及错误信息
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

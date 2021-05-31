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

// 请求返回响应结构体(swag)
type _response struct {
	Code    ResCode `json:"code" example:"400"`
	Message string  `json:"message" example:"status bad request"`
	Data    string  `json:"data,omitempty"`
}

// ResponseData 请求返回响应结构体
type ResponseData struct {
	Code    ResCode     `json:"code" example:"400"`
	Message interface{} `json:"message" example:"status bad request"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseSuccess 请求返回响应成功，返回成功信息
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	})
}

// ResponseError 请求返回响应错误，返回错误信息
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code:    code,
		Message: code.Msg(),
		Data:    nil,
	})
}

// ResponseErrorWithMsg 请求返回响应错误，返回错误及错误信息
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

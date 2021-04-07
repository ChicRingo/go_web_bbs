package controller

import "go_web_bbs/models"

// 帖子详情请求返回响应结构体
type _responsePostDetail struct {
	Code    ResCode               `json:"code" example:"400"`
	Message string                `json:"message" example:"status bad request"`
	Data    *models.ApiPostDetail `json:"data"`
}

// 帖子列表请求返回响应结构体
type _responsePostList struct {
	Code    ResCode                 `json:"code" example:"400"`
	Message string                  `json:"message" example:"status bad request"`
	Data    []*models.ApiPostDetail `json:"data"`
}

// 社区详情请求返回响应结构体
type _responseCommunityDetail struct {
	Code    ResCode           `json:"code" example:"400"`
	Message string            `json:"message" example:"status bad request"`
	Data    *models.Community `json:"data"`
}

// 社区列表请求返回响应结构体
type _responseCommunityList struct {
	Code    ResCode                   `json:"code" example:"400"`
	Message string                    `json:"message" example:"status bad request"`
	Data    []*models.CommunityDetail `json:"data"`
}

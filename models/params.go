package models

// 定义请求的参数

// ParamSingUp 注册请求参数
type ParamSingUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 投票数据
type ParamVoteData struct {
	//UserID 从请求中获取当前用户
	PostID    string `json:"post_id" binding:"required"`              // 帖子id，string防止前端处理失真
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票（1） 反对票（-1） 取消投票（0）
}

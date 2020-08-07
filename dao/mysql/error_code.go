package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorInvalidID       = errors.New("无效的ID")
	ErrorGenIDFailed     = errors.New("创建用户ID失败")
	ErrorQueryFailed     = errors.New("查询数据失败")
	ErrorInsertFailed    = errors.New("插入数据失败")
)

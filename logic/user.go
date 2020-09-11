package logic

import (
	"go_web_bbs/dao/mysql"
	"go_web_bbs/models"
	"go_web_bbs/pkg/jwt"
	"go_web_bbs/pkg/snowflake"
)

// 存放业务逻辑的代码

func SingUp(p *models.ParamSingUp) (err error) {

	// 1.判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2.生成UID
	userID := snowflake.GenID()
	// 构造user实例
	user := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3.保存进数据库
	return mysql.InsertUser(&user)
	// redis.xxx
}

// Login 判断用户密码是否正确
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到p.UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成JWT
	token, err := jwt.ReleaseToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}

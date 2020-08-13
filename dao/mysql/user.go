package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"go_web_bbs/models"
)

// 把每一步数据库操作封装成函数
// 待logic层更具业务需求调用

const secret = "secret"

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// InsertUser 新增用户插入一条新的纪录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// CheckUserPassword 检查指定用户名的密码是否正确
func CheckUserPassword(user *models.ParamLogin) error {
	//对登陆的密码进行加密
	password := encryptPassword(user.Password)
	sqlStr := `select password from user where username = ?`
	if err := db.Get(user, sqlStr, user.Username); err != nil {
		return err
	}
	if password != user.Password {
		return errors.New("密码不正确")
	}
	return nil
}

// Login 用户登录检查用户是否存在和密码是否正确
func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登录的密码
	sqlStr := `select username, password from user where username = ?`
	err = db.Get(user, sqlStr, user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return err
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserByID 根据用户id获取用户详细信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username
	from user
	where username = ?`
	err = db.Get(user, sqlStr, uid)
	//if err == sql.ErrNoRows {
	//	return ErrorUserNotExist
	//}
	//if err != nil {
	//	// 查询数据库失败
	//	return err
	//}
	return
}

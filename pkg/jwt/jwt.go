package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/spf13/viper"
)

var mySecret = []byte("你好骚啊")

/*
自定义声明结构体并内嵌jwt.StandardClaims
jwt包自带的jwt.StandardClaims只包含了官方字段
我们这里需要额外记录user_id,username字段，所以要自定义结构体
如果想要保存更多信息，都可以添加到这个结构体中
*/
type Claims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成token
func GenToken(userID int64, username string) (string, error) {
	// 创建一个我们自己的声明数据
	claims := Claims{
		userID,   // 自定义字段
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(viper.GetDuration("auth.jwt_expire") * time.Hour).Unix(), // 过期时间
			Issuer:    "go_web_bbs",                                                            // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// 解析token
func ParseToken(tokenString string) (*Claims, error) {
	var claims = new(Claims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

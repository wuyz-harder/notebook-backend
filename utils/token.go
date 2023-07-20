package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 自定义一个字符串
var jwtkey = []byte("www.topgoer.com")
var str string

type Claims struct {
	UserName string
	UserID   int
	jwt.StandardClaims
}

// 颁发token
func GenerateToken(name string, userID int) (string, error) {
	// 7天过期
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserName: name,
		UserID:   userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "127.0.0.1",  // 签名颁发者
			Subject:   "user token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// fmt.Println(token)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtkey, nil
	})
	fmt.Println("-====")
	fmt.Println(Claims.UserID)
	fmt.Println("-====")
	fmt.Println(Claims.UserName)
	fmt.Println("-====")
	return token, Claims, err
}

package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/phoenix-next/phoenix-server/global"
)

// GenerateToken 生成一个token
func GenerateToken(email string) (signedToken string) {
	claims := jwt.StandardClaims{
		Issuer:   "phoenix-server",
		Audience: email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := global.VP.GetString("server.secret")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		global.LOG.Panic("GenerateToken: sign token error")
	}
	return
}

//ValidateToken 验证token的正确性，正确则返回email
func ValidateToken(signedToken string) (email string, err error) {
	secret := global.VP.GetString("server.secret")
	token, err := jwt.Parse(
		signedToken,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil || !token.Valid {
		fmt.Println(err)
		err = errors.New("token isn't valid")
		return
	}
	email = token.Claims.(jwt.MapClaims)["aud"].(string)
	return
}

package utils

import (
	"errors"
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
	if err != nil {
		global.LOG.Panic("ValidateToken: parse token error")
	}
	if !token.Valid {
		err = errors.New("token isn't valid")
		return
	}
	email = token.Claims.(jwt.StandardClaims).Audience
	return
}

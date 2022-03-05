package utils

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/phoenix-next/phoenix-server/global"
	"strconv"
)

// GenerateToken 生成一个token
func GenerateToken(id uint64) (signedToken string) {
	claims := jwt.StandardClaims{
		Issuer:   "phoenix-server",
		Audience: strconv.FormatUint(id, 10),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := global.VP.GetString("server.secret")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		global.LOG.Panic("GenerateToken: sign token error")
	}
	return
}

//ValidateToken 验证token的正确性，正确则返回id
func ValidateToken(signedToken string) (id uint64, err error) {
	secret := global.VP.GetString("server.secret")
	token, err := jwt.Parse(
		signedToken,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil || !token.Valid {
		err = errors.New("token isn't valid")
		return
	}
	id, err = strconv.ParseUint(token.Claims.(jwt.MapClaims)["aud"].(string), 10, 64)
	if err != nil {
		err = errors.New("token isn't valid")
	}
	return
}

package middleware

import (
	"github.com/dgrijalva/jwt-go"
	token2 "ini/pkg/token"
	"time"
)

func GenerateToken(UID string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(300 * time.Hour)
	issuer := "KitZhangYs"
	claims := token2.MyCustomClaims{
		Username: UID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("OYHX"))
	return token, err
}

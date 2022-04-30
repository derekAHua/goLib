package utils

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

const singKey = "111"

func TestJWT_CreateToken(t *testing.T) {
	j := NewJWT(singKey)

	token, err := j.CreateToken(CustomClaims{
		ID:          1,
		NickName:    "1",
		AuthorityId: 1,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // 签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30天过期
			Issuer:    "YZH",                           // 哪个机构进行的签名
			IssuedAt:  time.Now().Unix(),               // 发放时间
		},
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(token)

	parseToken, err := j.ParseToken(token)
	if err != nil {
		t.Error(err)
	}
	t.Log(parseToken.String())
}

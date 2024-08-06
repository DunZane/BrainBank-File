package utils

import (
	"github.com/dunzane/brainbank-file/api/internal/config"
	"reflect"
	"testing"
)

func TestJWT(t *testing.T) {
	userId := 873526635
	email := "test_email@brainbank.com"

	c := config.Config{
		Auth: struct {
			AccessSecret string
			AccessExpire int64
		}{
			AccessSecret: "my-secret-for-brainbank",
			AccessExpire: 3600,
		},
	}

	token, err := GenerateToken(c, int64(userId), email)
	if err != nil {
		t.Errorf("generate token failed:%v", err)
	} else {
		t.Logf("generate token success:%s", token)
	}

	claims, err := ParseToken(c, token)
	if err != nil {
		t.Errorf("parse token failed:%v", err)
	}

	if !reflect.DeepEqual(claims.Email, email) || !reflect.DeepEqual(claims.UserId, int64(userId)) {
		t.Errorf("there are some error logic in jwt")
	} else {
		t.Log("jwt test success.")
	}
}

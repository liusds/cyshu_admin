package service

import (
	"context"
	"crypto/md5"
	v1 "cyshu_admin/api/admin/v1"
	"cyshu_admin/app/cyshu/internal/biz"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string
	Password string
	jwt.StandardClaims
}

func (s *CyshuAdmin) Register(ctx context.Context, in *v1.RegisterRequest) (*v1.RegisterReply, error) {
	respond, err := s.addUser(ctx, "v1", in.User)
	if err != nil {
		return nil, err
	}
	return &v1.RegisterReply{Respond: respond}, nil
}

func (s *CyshuAdmin) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginReply, error) {

	var user *biz.User
	if len(in.Username) > 0 && len(in.Password) >= 6 {
		user = &biz.User{
			UserName: in.Username,
			Password: in.Password,
		}
	} else {
		result := &v1.Result{
			Code:    http.StatusBadRequest,
			Token:   "",
			Message: "账号或密码不能为空",
		}
		return &v1.LoginReply{Result: result}, nil
	}
	if err := s.au.Login(ctx, user); err != nil {
		return nil, err
	}

	// Token .
	claims := Claims{
		encodeMD5(user.UserName),
		encodeMD5(user.Password),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    strconv.Itoa(int(user.ID)),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(SecretKey))
	if err != nil {
		return nil, v1.ErrorApiTokenError("could not login")
	}

	// cookie := http.Cookie{
	// 	Name:     user.UserName,
	// 	Value:    token,
	// 	Expires:  time.Now().Add(time.Hour * 24),
	// 	HttpOnly: true,
	// }

	// fmt.Println(cookie)

	result := &v1.Result{
		Code:    http.StatusOK,
		Token:   token,
		Message: "会员登录成功",
	}

	return &v1.LoginReply{Result: result}, nil
}

func (s *CyshuAdmin) Logout(ctx context.Context, in *v1.LogoutRequest) (*v1.LogoutReply, error) {
	if err := s.au.Logout(ctx, in.GetId()); err != nil {
		return nil, err
	}
	respond := &v1.Respond{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "会员退出成功",
	}
	return &v1.LogoutReply{
		Api:     apiVersion,
		Respond: respond,
	}, nil
}

func encodeMD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	md := m.Sum(nil)
	return hex.EncodeToString(md)
}

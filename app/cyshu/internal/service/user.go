package service

import (
	"context"
	v1 "cyshu_admin/api/admin/v1"
	"cyshu_admin/app/cyshu/internal/biz"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	apiVersion = "v1"
	SecretKey  = "cyshu_v1"
)

func (s *CyshuAdmin) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return v1.ErrorApiVersionNotUnsupported("unsupported api version")
		}
	} else {
		return v1.ErrorApiVersionNotUnsupported("unsupported api version")
	}
	return nil
}

func (s *CyshuAdmin) addUser(ctx context.Context, api string, u *v1.User) (*v1.Respond, error) {
	if err := s.checkAPI(api); err != nil {
		return nil, err
	}
	var user *biz.User
	isEmail, err := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, u.Username)
	if err != nil {
		return nil, err
	}
	if isEmail {
		if u.Username == "" || u.Password == "" {
			respond := &v1.Respond{
				Code:    http.StatusBadRequest,
				Status:  "bad request",
				Message: "会员名称或密码不能为空",
			}
			return respond, nil
		} else if len(u.Username) < 3 || len(u.Password) < 6 {
			respond := &v1.Respond{
				Code:    http.StatusBadRequest,
				Status:  "bad request",
				Message: "账号长度大于3位，密码长度大于6位",
			}
			return respond, nil
		} else {
			pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost) //加密处理
			if err != nil {
				return nil, err
			}
			user = &biz.User{
				UserName:   u.Username,
				Password:   string(pwd),
				Role:       u.Role,
				UserPhoto:  u.UserPhoto,
				UserPhone:  int(u.UserPhone),
				CreateTime: time.Now().Unix(),
			}
		}
	} else {
		respond := &v1.Respond{
			Code:    http.StatusBadRequest,
			Status:  "bad request",
			Message: "请用邮箱注册",
		}
		return respond, nil
	}

	res, err := s.uc.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	respond := &v1.Respond{
		Code:    http.StatusOK,
		Status:  strconv.Itoa(res),
		Message: "会员注册成功",
	}
	return respond, nil
}

func (s *CyshuAdmin) CreateUsers(ctx context.Context, in *v1.CreateUsersRequest) (*v1.CreateUsersReply, error) {
	respond, err := s.addUser(ctx, in.Api, in.User)
	if err != nil {
		return nil, err
	}
	return &v1.CreateUsersReply{Api: apiVersion, Respond: respond}, nil
}

func (s *CyshuAdmin) GetUsers(ctx context.Context, in *v1.GetUsersRequest) (*v1.GetUsersReply, error) {
	if err := s.checkAPI(in.Api); err != nil {
		return nil, err
	}
	u, err := s.uc.Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	user := &v1.User{
		Id:        u.ID,
		Username:  u.UserName,
		Role:      u.Role,
		UserPhoto: u.UserPhoto,
		UserPhone: int64(u.UserPhone),
	}
	return &v1.GetUsersReply{Api: apiVersion, User: user}, nil
}

func (s *CyshuAdmin) UpdateUsers(ctx context.Context, in *v1.UpdateUsersRequest) (*v1.UpdateUsersReply, error) {
	if err := s.checkAPI(in.Api); err != nil {
		return nil, err
	}
	var user *biz.User
	isEmail, err := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, in.User.Username)
	if err != nil {
		return nil, err
	}
	if len(in.User.Username) <= 0 {
		return nil, v1.ErrorUserNotUpdate("账号不能为空")
	} else if !isEmail {
		return nil, v1.ErrorUserNotUpdate("请填写邮箱")
	} else {
		user = &biz.User{
			ID:        in.User.GetId(),
			UserName:  in.User.Username,
			Role:      in.User.Role,
			UserPhoto: in.User.UserPhoto,
			UserPhone: int(in.User.UserPhone),
		}
	}
	if err := s.uc.Update(ctx, user); err != nil {
		return nil, err
	}

	respond := &v1.Respond{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "会员更新成功",
	}
	return &v1.UpdateUsersReply{Api: apiVersion, Respond: respond}, nil
}

func (s *CyshuAdmin) DeleteUsers(ctx context.Context, in *v1.DeleteUsersRequest) (*v1.DeleteUsersReply, error) {
	if err := s.checkAPI(in.Api); err != nil {
		return nil, err
	}
	if err := s.uc.Delete(ctx, in.GetId()); err != nil {
		return nil, err
	}

	respond := &v1.Respond{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "会员删除成功",
	}
	return &v1.DeleteUsersReply{Api: apiVersion, Respond: respond}, nil
}

func (s *CyshuAdmin) ListUsers(ctx context.Context, in *v1.ListUsersRequest) (*v1.ListUsersReply, error) {
	if err := s.checkAPI(in.Api); err != nil {
		return nil, err
	}
	users, err := s.uc.List(ctx)
	if err != nil {
		return nil, err
	}
	listUsers := []*v1.User{}
	for _, v := range users {
		listUsers = append(listUsers, &v1.User{Id: v.ID, Username: v.UserName, Password: "", Role: v.Role, UserPhoto: v.UserPhoto, UserPhone: int64(v.UserPhone)})
	}
	return &v1.ListUsersReply{Api: apiVersion, User: listUsers}, nil
}

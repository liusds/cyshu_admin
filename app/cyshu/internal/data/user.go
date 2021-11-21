package data

import (
	"context"
	v1 "cyshu_admin/api/admin/v1"
	"cyshu_admin/app/cyshu/internal/biz"
	"encoding/json"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) CreateUser(ctx context.Context, u *biz.User) (int, error) {
	if r.data.db.Where("user_name = ?", u.UserName).Find(&biz.User{}).RowsAffected > 0 {
		return 1, v1.ErrorUserNotCreate("该会员已存在")
	}
	result := r.data.db.Model(&biz.User{}).Create(&u).RowsAffected
	if result == 0 {
		return 0, v1.ErrorUserNotCreate("会员添加失败")
	}
	return int(result), nil
}

func (r *userRepo) UpdateUser(ctx context.Context, u *biz.User) error {
	i, err := r.data.redisClient.Exists(u.UserName).Result()
	if err != nil {
		return err
	}
	if i == 0 {
		userData, err := json.Marshal(u)
		if err != nil {
			return err
		}
		_, err = r.data.redisClient.Set(u.UserName, userData, time.Second*300).Result()
		if err != nil {
			return err
		}
	} else {
		user := &biz.User{}
		resBytes, err := r.data.redisClient.Get(u.UserName).Bytes()
		if err != nil {
			return err
		}
		if err := json.Unmarshal(resBytes, user); err != nil {
			return err
		}
		if user.ID == u.ID && user.UserName == u.UserName && user.Role == u.Role && user.UserPhoto == u.UserPhoto && user.UserPhone == u.UserPhone {
			return v1.ErrorUserNotUpdate("请勿重复修改")
		} else {
			userData, err := json.Marshal(u)
			if err != nil {
				return err
			}
			_, err = r.data.redisClient.Set(u.UserName, userData, time.Second*10).Result()
			if err != nil {
				return err
			}
		}
	}
	result := r.data.db.Model(&biz.User{}).Where("id = ?", u.ID).Updates(&u).RowsAffected
	if result == 0 {
		return v1.ErrorUserNotUpdate("会员修改失败")
	}
	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id int64) error {
	result := r.data.db.Where("id = ?", id).Delete(&biz.User{}).RowsAffected
	if result == 0 {
		return v1.ErrorUserNotDelete("会员删除失败")
	}
	return nil
}

func (r *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	var user *biz.User
	result := r.data.db.Table("users").Where("id = ?", id).Scan(&user).RowsAffected
	if result == 0 {
		return nil, v1.ErrorUserNotFound("没有该会员")
	}
	return user, nil
}

func (r *userRepo) ListUser(ctx context.Context) ([]*biz.User, error) {
	userList := []*biz.User{}
	r.data.db.Table("users").Find(&userList)
	return userList, nil
}

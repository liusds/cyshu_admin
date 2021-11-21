package data

import (
	"context"
	v1 "cyshu_admin/api/admin/v1"
	"cyshu_admin/app/cyshu/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
)

type adminRepo struct {
	data *Data
	log  *log.Helper
}

func NewAdminRepo(data *Data, logger log.Logger) biz.AdminRepo {
	return &adminRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *adminRepo) Login(ctx context.Context, u *biz.User) error {
	var user *biz.User
	result := r.data.db.Table("users").Where("user_name = ?", u.UserName).Find(&user).RowsAffected
	if result == 0 {
		return v1.ErrorUserNotFound("会员账号错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return v1.ErrorUserPasswordError("密码错误")
	}
	return nil
}

func (r *adminRepo) Logout(ctx context.Context, id int64) error {
	return nil
}

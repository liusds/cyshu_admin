package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type AdminRepo interface {
	Login(context.Context, *User) error
	Logout(context.Context, int64) error
}

type AdminUsecase struct {
	repo AdminRepo
	log  *log.Helper
}

func NewAdminUsecase(repo AdminRepo, logger log.Logger) *AdminUsecase {
	return &AdminUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *AdminUsecase) Login(ctx context.Context, u *User) error {
	return uc.repo.Login(ctx, u)
}

func (uc *AdminUsecase) Logout(ctx context.Context, id int64) error {
	return uc.repo.Logout(ctx, id)
}

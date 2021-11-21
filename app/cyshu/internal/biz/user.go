package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID         int64
	UserName   string
	Password   string
	Role       string
	UserPhoto  string
	Status     int
	UserPhone  int
	LastIP     string
	LastTime   int64
	CreateTime int64
}

type UserRepo interface {
	CreateUser(context.Context, *User) (int, error)
	UpdateUser(context.Context, *User) error
	DeleteUser(context.Context, int64) error
	GetUser(context.Context, int64) (*User, error)
	ListUser(context.Context) ([]*User, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) Create(ctx context.Context, u *User) (int, error) {
	return uc.repo.CreateUser(ctx, u)
}

func (uc *UserUsecase) Update(ctx context.Context, u *User) error {
	return uc.repo.UpdateUser(ctx, u)
}

func (uc *UserUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteUser(ctx, id)
}

func (uc *UserUsecase) Get(ctx context.Context, id int64) (*User, error) {
	return uc.repo.GetUser(ctx, id)
}

func (uc *UserUsecase) List(ctx context.Context) ([]*User, error) {
	return uc.repo.ListUser(ctx)
}

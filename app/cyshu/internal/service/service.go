package service

import (
	v1 "cyshu_admin/api/admin/v1"
	"cyshu_admin/app/cyshu/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewCyshuAdminService)

type CyshuAdmin struct {
	v1.UnimplementedCyshuAdminServer
	log *log.Helper
	uc  *biz.UserUsecase
	au  *biz.AdminUsecase
}

func NewCyshuAdminService(uc *biz.UserUsecase, au *biz.AdminUsecase, logger log.Logger) *CyshuAdmin {
	return &CyshuAdmin{
		uc:  uc,
		au:  au,
		log: log.NewHelper(logger),
	}
}

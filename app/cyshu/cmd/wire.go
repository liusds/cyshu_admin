// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"cyshu_admin/app/cyshu/internal/biz"
	"cyshu_admin/app/cyshu/internal/conf"
	"cyshu_admin/app/cyshu/internal/data"
	"cyshu_admin/app/cyshu/internal/server"
	"cyshu_admin/app/cyshu/internal/service"

	// "cyshu_admin_users/log"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

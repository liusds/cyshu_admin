package server

import (
	v1 "cyshu_admin/api/admin/v1"
	"cyshu_admin/app/cyshu/internal/conf"
	"cyshu_admin/app/cyshu/internal/middleware/jwt"
	"cyshu_admin/app/cyshu/internal/service"
	"regexp"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, s *service.CyshuAdmin, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			metadata.Server(),
			selector.Server(
				jwt.Server(),
			).Match(func(operation string) bool {
				r, err := regexp.Compile(OPERATION)
				if err != nil {
					return true
				}
				return r.FindString(operation) != operation
			}).Build(),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	v1.RegisterCyshuAdminHTTPServer(srv, s)
	return srv
}

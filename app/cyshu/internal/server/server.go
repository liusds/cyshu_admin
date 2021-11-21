package server

import (
	"github.com/google/wire"
)

const OPERATION = `/api.admin.v1.CyshuAdmin/Login|/api.admin.v1.CyshuAdmin/Register`

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)

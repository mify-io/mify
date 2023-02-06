package router

import (
	"net/http"

	"{{.CoreInclude}}"
)


type routerConfig struct {
	Middlewares []func(http.Handler) http.Handler
}

func NewRouterConfig(ctx *core.MifyServiceContext) *routerConfig {
	return &routerConfig {
		Middlewares: []func(http.Handler) http.Handler {
		// Add your middlewares here
		},
	}
}

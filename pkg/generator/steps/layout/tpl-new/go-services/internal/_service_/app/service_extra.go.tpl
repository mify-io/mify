// vim: set ft=go:

package app

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

type ServiceExtra struct {
	// Append your dependencies here
}

func NewServiceExtra(ctx *core.MifyServiceContext) (*ServiceExtra, error) {
	// Here you can do your custom service initialization, prepare dependencies
	extra := &ServiceExtra{
		// Here you can initialize your dependencies
	}
	return extra, nil
}

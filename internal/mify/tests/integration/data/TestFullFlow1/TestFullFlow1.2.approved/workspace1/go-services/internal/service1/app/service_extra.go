package app

import (
	"net/http"

	"example.com/namespace/workspace1/go-services/internal/service1/generated/core"
)

type routerConfig struct {
	Middlewares []func(http.Handler) http.Handler
}

func NewRouterConfig() *routerConfig {
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

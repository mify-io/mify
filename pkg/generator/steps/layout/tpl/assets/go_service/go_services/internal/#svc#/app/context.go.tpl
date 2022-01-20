package app

import "net/http"

var (
	MiddlewareList []func(http.Handler) http.Handler
)

type ServiceContext struct {
	// Append your dependencies here
}

func NewServiceContext() (*ServiceContext, error) {
	// Here you can do your custom service initialization, prepare dependencies, middlewares
	// Add middlewares to MiddlewareList
	context := &ServiceContext{
		// Here you can initialize your dependencies
	}
	return context, nil
}

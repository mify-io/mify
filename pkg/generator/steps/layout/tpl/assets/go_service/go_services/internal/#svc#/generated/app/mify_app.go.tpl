{{- .Workspace.TplHeader}}

package app

import (
	"net/http"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"{{.GoModule}}/internal/{{.ServiceName}}/app"
	"{{.GoModule}}/internal/{{.ServiceName}}/generated/api/init"
	"{{.GoModule}}/internal/{{.ServiceName}}/generated/core"
)

type routerConfig struct {
	middlewares []func(http.Handler) http.Handler
}

func newRouterConfig() *routerConfig {
	conf := app.NewRouterConfig()
	return &routerConfig {
		middlewares: conf.Middlewares,
	}
}

func (r *routerConfig) Middlewares() []func(http.Handler) http.Handler {
	return r.middlewares
}

type MifyServiceApp struct {
	context *core.MifyServiceContext
	router  chi.Router
}

func NewMifyServiceApp(goGontext context.Context) *MifyServiceApp {
	serviceContext, _ := core.NewMifyServiceContext(
		goGontext, "{{.ServiceName}}",
		func(ctx *core.MifyServiceContext) (interface{}, error) {
			return app.NewServiceContext(ctx)
		})
	router := openapi_init.Routes(serviceContext, newRouterConfig())

	router.Handle("/metrics", promhttp.Handler())

	return &MifyServiceApp{
		context: serviceContext,
		router:  router,
	}
}

func (app MifyServiceApp) Run() {
	app.context.Logger().Info("Starting...")

	err := openapi_init.RunServer(app.context, app.router)
	if err != nil {
		app.context.Logger().Panic("failed to run server", zap.Error(err))
	}
}

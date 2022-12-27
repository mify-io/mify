// vim: set ft=go:

{{- .TplHeader}}
// vim: set ft=go:

package app

import (
	"net/http"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"{{.AppImportPath}}"
	"{{.InitImportPath}}"
	"{{.CoreImportPath}}"
	"{{.ApiImportPath}}"
)

type routerConfig struct {
	middlewares []func(http.Handler) http.Handler
}

func newRouterConfig(ctx *core.MifyServiceContext) *routerConfig {
	conf := app.NewRouterConfig(ctx)
	return &routerConfig {
		middlewares: conf.Middlewares,
	}
}

func (r *routerConfig) Middlewares() []func(http.Handler) http.Handler {
	return r.middlewares
}

type MifyServiceApp struct {
	context *core.MifyServiceContext
	maintenanceRouter chi.Router
	apiRouter  chi.Router
}

type maintenanceRouter struct {}

func (r maintenanceRouter) Routes() openapi.Routes {
	return []openapi.Route {
		{
			Name: "metrics",
			Method: "post",
			Pattern: "/metrics",
			HandlerFunc: func(rw http.ResponseWriter, r *http.Request) {
				promhttp.Handler().ServeHTTP(rw, r)
			},
		},
		{
			Name: "swagger-ui",
			Method: "get",
			Pattern: "/swagger-ui/*",
			HandlerFunc: openapi.SwaggerUIHandlerFunc,
		},
	}
}

func NewMifyServiceApp(goGontext context.Context) *MifyServiceApp {
	serviceContext, err := core.NewMifyServiceContext(
		goGontext, "{{.ServiceName}}",
		func(ctx *core.MifyServiceContext) (interface{}, error) {
			return app.NewServiceExtra(ctx)
		})
	if err != nil {
		panic(err)
	}

	return &MifyServiceApp{
		context: serviceContext,
		maintenanceRouter: openapi.NewRouter(serviceContext, newRouterConfig(serviceContext), maintenanceRouter{}),
		apiRouter: openapi_init.Routes(serviceContext, newRouterConfig(serviceContext)),
	}
}

func (app MifyServiceApp) Run() {
	app.context.Logger().Info("Starting...")

	ctx, cancel := context.WithCancel(app.context.GetContext())

	go func() {
		err := runMaintenanceServer(app.context, app.maintenanceRouter)
		if err != nil {
			app.context.Logger().Panic("failed to run maintenance server", zap.Error(err))
		}

		cancel()
	}()

	go func() {
		err := runApiServer(app.context, app.apiRouter)
		if err != nil {
			app.context.Logger().Panic("failed to run api server", zap.Error(err))
		}

		cancel()
	}()

	<-ctx.Done()
}

{{- .Workspace.TplHeader}}

package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"{{.GoModule}}/internal/{{.ServiceName}}/app"
	"{{.GoModule}}/internal/{{.ServiceName}}/generated/api/init"
	"{{.GoModule}}/internal/{{.ServiceName}}/generated/core"
)

type MifyServiceApp struct {
	context *core.MifyServiceContext
	router  chi.Router
}

func NewMifyServiceApp(goGontext context.Context) *MifyServiceApp {
	serviceContext, _ := core.NewMifyServiceContext(goGontext, "{{.ServiceName}}")
	router := openapi_init.Routes(serviceContext)

	for _, middleware := range app.MiddlewareList {
		router.Use(middleware)
	}
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

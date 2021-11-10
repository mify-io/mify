{{- .Workspace.TplHeader}}

package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"{{.GoModule}}/internal/{{.ServiceName}}/generated/api/init"
	"{{.GoModule}}/internal/{{.ServiceName}}/generated/core"
)

type MifyServiceApp struct {
	context *core.MifyServiceContext
	router  chi.Router
}

func NewMifyServiceApp() *MifyServiceApp {
	serviceContext, _ := core.NewMifyServiceContext("{{.ServiceName}}")
	router := openapi_init.Routes(serviceContext)

	router.Handle("/metrics", promhttp.Handler())

	return &MifyServiceApp{
		context: serviceContext,
		router:  router,
	}
}

func (app MifyServiceApp) Run() {
	app.context.Logger().Info("Starting...")

	err := http.ListenAndServe(":8080", app.router)
	if err != nil {
		app.context.Logger().Panic("failed to listen", zap.Error(err))
	}
}

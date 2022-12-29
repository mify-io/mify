{{- .TplHeader}}
// vim: set ft=go:

package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"{{.CoreImportPath}}"
	"{{.ConfigsImportPath}}"
)

type ServerConf struct {
	ServerEndpoint string `envconfig:"{{.ApiEndpointEnv}}" default:"{{.ApiEndpoint}}"`
	MaintenanceEndpoint string `envconfig:"{{.MaintenanceEndpointEnv}}" default:"{{.MaintenanceEndpoint}}"`
}

func getServerConf(conf *configs.MifyStaticConfig) *ServerConf {
	return conf.MustGetPtr((*ServerConf)(nil)).(*ServerConf)
}

func runApiServer(ctx *core.MifyServiceContext, router chi.Router) error {
	conf := getServerConf(ctx.StaticConfig())
	return runServer(ctx, "api", conf.ServerEndpoint, router)
}

func runMaintenanceServer(ctx *core.MifyServiceContext, router chi.Router) error {
	conf := getServerConf(ctx.StaticConfig())
	return runServer(ctx, "maintenance", conf.MaintenanceEndpoint, router)
}

func runServer(ctx *core.MifyServiceContext, name string, endpoint string, router chi.Router) error {
	ctx.Logger().Info(fmt.Sprintf("starting %s server", name), zap.String("endpoint", endpoint))
	err := http.ListenAndServe(endpoint, router)
	if err != nil {
		ctx.Logger().Error("failed to listen", zap.Error(err))
		return err
	}
	return nil
}

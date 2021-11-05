{{- .Workspace.TplHeader}}

package main

import (
	"net/http"

	"go.uber.org/zap"
	"{{.Repository}}/go_services/internal/{{.ServiceName}}/generated/core"
	"repo.com/namespace/somerepo/go_services/internal/service1/generated/api/init"
)

func main() {
	serviceContext, _ := core.NewMifyServiceContext("service1")
	serviceContext.Logger().Info("Starting...")

	// tmp
	router := openapi_init.Routes(serviceContext)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		serviceContext.Logger().Panic("failed to listen", zap.Error(err))
	}
}

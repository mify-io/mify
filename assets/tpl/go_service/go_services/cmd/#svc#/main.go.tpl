{{- .Workspace.TplHeader}}

package main

import (
	"{{.Repository}}/go_services/internal/{{.ServiceName}}/generated/core"
)

func main() {
	serviceContext, _ := core.NewMifyServiceContext("service1")
	serviceContext.Logger.Info("Starting...")

	requestContext, _ := core.NewMifyRequestContext(serviceContext)
	requestContext.Logger.Info("Test")
}

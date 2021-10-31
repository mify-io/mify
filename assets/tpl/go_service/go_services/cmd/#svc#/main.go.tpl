package main

import (
	"{{.Repository}}/go_services/internal/{{.ServiceName}}/misc"

	"go.uber.org/zap"
)

func main() {
	context, _ := misc.NewMifyServiceContext("{{.ServiceName}}")

	context.Logger.Info("Starting...", zap.String("service_name", context.ServiceName))
}

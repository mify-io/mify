package main

import (
	"internal/{{.ServiceName}}/misc"

	"go.uber.org/zap"
)

func main() {
	context := misc.InitContext("{{.ServiceName}}")

	context.Logger.Info("Starting...", zap.String("service_name", context.ServiceName))
}

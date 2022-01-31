package core

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type serviceContextModel struct {
	TplHeader          string
	ConfigsImportPath  string
	LogsImportPath     string
	MetricsImportPath  string
	ConsulImportPath   string
	PostgresImportPath string
}

func newServiceContextModel(ctx *gencontext.GenContext) serviceContextModel {
	// TODO: move paths to description

	postgresImportPath := ""
	if ctx.GetMifySchema().Postgres.Enabled {
		postgresImportPath = fmt.Sprintf("%s/internal/pkg/generated/postgres", ctx.GetWorkspace().GetGoModule())
	}

	return serviceContextModel{
		TplHeader:          ctx.GetWorkspace().TplHeader,
		ConfigsImportPath:  fmt.Sprintf("%s/internal/pkg/generated/configs", ctx.GetWorkspace().GetGoModule()),
		LogsImportPath:     fmt.Sprintf("%s/internal/pkg/generated/logs", ctx.GetWorkspace().GetGoModule()),
		MetricsImportPath:  fmt.Sprintf("%s/internal/pkg/generated/metrics", ctx.GetWorkspace().GetGoModule()),
		ConsulImportPath:   fmt.Sprintf("%s/internal/pkg/generated/consul", ctx.GetWorkspace().GetGoModule()),
		PostgresImportPath: postgresImportPath,
	}
}

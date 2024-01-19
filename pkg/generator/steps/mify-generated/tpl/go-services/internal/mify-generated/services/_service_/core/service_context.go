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
		servicePath := ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetServicePackage()
		postgresImportPath = fmt.Sprintf("%s/postgres", servicePath)
	}
	commonPath := ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetCommonPackage()

	return serviceContextModel{
		TplHeader:          ctx.GetWorkspace().TplHeader,
		ConfigsImportPath:  fmt.Sprintf("%s/configs", commonPath),
		LogsImportPath:     fmt.Sprintf("%s/logs", commonPath),
		MetricsImportPath:  fmt.Sprintf("%s/metrics", commonPath),
		ConsulImportPath:   fmt.Sprintf("%s/consul", commonPath),
		PostgresImportPath: postgresImportPath,
	}
}

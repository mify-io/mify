{{- .TplHeader}}
// vim: set ft=go:

package core

import (
	"os"
	"context"

	"go.uber.org/zap"
	"github.com/hashicorp/consul/api"
	"{{.ConfigsImportPath}}"
	"{{.LogsImportPath}}"
	"{{.MetricsImportPath}}"
	"{{.ConsulImportPath}}"
{{- if .PostgresImportPath}}
	"github.com/jackc/pgx/v4/pgxpool"
	"{{.PostgresImportPath}}"
{{- end}}
)

type MifyServiceContext struct {
	goContext   context.Context
	serviceName string
	hostname    string

	loggerWrapper  *logs.MifyLoggerWrapper
	metricsWrapper *metrics.MifyMetricsWrapper
	staticConfig   *configs.MifyStaticConfig
	dynamicConfig  *configs.MifyDynamicConfig
	clients        *MifyServiceClients
	{{if .PostgresImportPath}}postgres       *pgxpool.Pool{{end}}

	serviceExtra interface{}
}

type ServiceExtraCreateFunc = func(ctx *MifyServiceContext) (interface{}, error)

func NewMifyServiceContext(goContext context.Context, serviceName string, extraCreate ServiceExtraCreateFunc) (*MifyServiceContext, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return &MifyServiceContext{}, err
	}

	context := &MifyServiceContext{
		goContext:   goContext,
		serviceName: serviceName,
		hostname:    hostname,
	}

	staticConfig, err := configs.NewMifyStaticConfig()
	if err != nil {
		return nil, err
	}
	context.staticConfig = staticConfig

	logger, err := logs.NewMifyLoggerWrapper(context)
	if err != nil {
		return nil, err
	}
	context.loggerWrapper = logger

	consulClient, err := api.NewClient(&api.Config{Address: consul.GetConsulConfig(staticConfig).Endpoint})
	if err != nil {
		return nil, err
	}
	dynamicConfig, err := configs.NewMifyDynamicConfig(consulClient)
	if err != nil {
		return nil, err
	}
	context.dynamicConfig = dynamicConfig

	clients, err := NewMifyServiceClients(context)
	if err != nil {
		return nil, err
	}
	context.clients = clients

{{if .PostgresImportPath}}
	pgConnString := postgres.GetPostgresConfig(staticConfig).DatabaseUrl
	dbpool, err := pgxpool.Connect(goContext, pgConnString)
	if err != nil {
		return nil, err
	}
	context.postgres = dbpool
{{end}}

	extra, err := extraCreate(context)
	if err != nil {
		return nil, err
	}
	context.serviceExtra = extra

	return context, nil
}

func (c *MifyServiceContext) GetContext() context.Context {
	return c.goContext
}

func (c *MifyServiceContext) ServiceName() string {
	return c.serviceName
}

func (c *MifyServiceContext) Hostname() string {
	return c.hostname
}

func (c *MifyServiceContext) Logger() *zap.Logger {
	return c.loggerWrapper.Logger()
}

func (c *MifyServiceContext) LoggerFor(component string) *zap.Logger {
	return c.loggerWrapper.LoggerFor(component)
}

func (c *MifyServiceContext) StaticConfig() *configs.MifyStaticConfig {
	return c.staticConfig
}

func (c *MifyServiceContext) DynamicConfig() *configs.MifyDynamicConfig {
	return c.dynamicConfig
}

func (c *MifyServiceContext) Clients() *MifyServiceClients {
	return c.clients
}

{{- if .PostgresImportPath}}

func (c *MifyServiceContext) Postgres() *pgxpool.Pool {
	return c.postgres
}
{{- end}}

func (c *MifyServiceContext) ServiceExtra() interface{} {
	return c.serviceExtra
}

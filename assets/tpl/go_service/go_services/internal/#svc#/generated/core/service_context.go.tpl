{{- .Workspace.TplHeader}}

package core

import (
	"os"

	"go.uber.org/zap"
	"github.com/hashicorp/consul/api"
	"{{.GoModule}}/internal/pkg/generated/configs"
	"{{.GoModule}}/internal/pkg/generated/logs"
	"{{.GoModule}}/internal/{{.ServiceName}}/app"
)

type MifyServiceContext struct {
	serviceName string
	hostname    string

	loggerWrapper  *logs.MifyLoggerWrapper
	metricsWrapper *MifyMetricsWrapper
	staticConfig   *configs.MifyStaticConfig
	dynamicConfig  *configs.MifyDynamicConfig
	clients        *MifyServiceClients

	serviceContext *app.ServiceContext
}

func NewMifyServiceContext(serviceName string) (*MifyServiceContext, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return &MifyServiceContext{}, err
	}

	context := &MifyServiceContext{
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

	consulClient, err := api.NewClient(&api.Config{Address: GetConsulConfig(staticConfig).Endpoint})
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

	svcCtx, err := app.NewServiceContext()
	if err != nil {
		return nil, err
	}
	context.serviceContext = svcCtx

	return context, nil
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

func (c *MifyServiceContext) ServiceContext() *app.ServiceContext {
	return c.serviceContext
}

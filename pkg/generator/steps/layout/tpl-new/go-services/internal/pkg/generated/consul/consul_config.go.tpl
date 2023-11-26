{{- .TplHeader}}
// vim: set ft=go:

package consul

import "{{.Workspace.PackageName}}/internal/pkg/generated/configs"

type ConsulConfig struct {
	Endpoint string `yaml:"endpoint" envconfig:"CONSUL_ENDPOINT" default:"127.0.0.1:8500"`
}

func GetConsulConfig(cfg *configs.MifyStaticConfig) *ConsulConfig {
	return cfg.MustGetPtr((*ConsulConfig)(nil)).(*ConsulConfig)
}

package	config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	serviceConfigName = "service.mify.yaml"
)

type ServiceConfig struct {
	ServiceName string   `yaml:"service_name"`
	Maintainers []string `yaml:"maintainers"`
}

func ReadServiceConfig(workspaceDir string, serviceName string) (ServiceConfig, error) {
	// TODO: service index file?
	// need to find service yaml
	return ServiceConfig{}, nil
}

func SaveServiceConfig(path string, conf ServiceConfig) error {
	data, err := yaml.Marshal(&conf)
	if err != nil {
		return fmt.Errorf("failed to create service config: %w", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", path, serviceConfigName), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create service config: %w", err)
	}
	return nil
}

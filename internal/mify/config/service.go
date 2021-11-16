package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	serviceConfigName = "service.mify.yaml"
)

type ServiceOpenAPIClientConfig struct {}

type ServiceOpenAPIConfig struct {
	Clients map[string]ServiceOpenAPIClientConfig `yaml:"clients,omitempty"`
}

type ServiceConfig struct {
	ServiceName string   `yaml:"service_name"`
	Maintainers []string `yaml:"maintainers"`

	OpenAPI ServiceOpenAPIConfig `yaml:"openapi,omitempty"`
}

func ReadServiceConfig(workspaceDir string, serviceName string) (ServiceConfig, error) {
	goServicesConfigFile := filepath.Join(workspaceDir, "go_services/cmd/", serviceName, serviceConfigName)
	if _, err := os.Stat(goServicesConfigFile); errors.Is(err, os.ErrNotExist) {
		return ServiceConfig{}, fmt.Errorf("service %s not found", serviceName)
	}

	rawData, err := ioutil.ReadFile(goServicesConfigFile)
	if err != nil {
		return ServiceConfig{}, fmt.Errorf("failed to read service config: %w", err)
	}

	var data ServiceConfig
	err = yaml.Unmarshal(rawData, &data)
	if err != nil {
		return ServiceConfig{}, fmt.Errorf("failed to read service config: %w", err)
	}

	return data, nil
}

func SaveServiceConfig(workspaceDir string, serviceName string, conf ServiceConfig) error {
	path := filepath.Join(workspaceDir, "go_services/cmd", serviceName)
	data, err := yaml.Marshal(&conf)
	if err != nil {
		return fmt.Errorf("failed to create service config: %w", err)
	}

	err = ioutil.WriteFile(filepath.Join(path, serviceConfigName), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create service config: %w", err)
	}
	return nil
}

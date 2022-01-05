package mifyconfig

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ServiceOpenAPIClientConfig struct {}

type ServiceOpenAPIConfig struct {
	Clients map[string]ServiceOpenAPIClientConfig `yaml:"clients,omitempty"`
}

type ServiceConfig struct {
	Language ServiceLanguage `yaml:"-"`

	ServiceName string   `yaml:"service_name"`
	Maintainers []string `yaml:"maintainers"`

	OpenAPI ServiceOpenAPIConfig `yaml:"openapi,omitempty"`
}

func ReadServiceConfig(workspaceDir string, serviceName string) (ServiceConfig, error) {
	w, err := NewWorkspace(workspaceDir)
	if err != nil {
		return ServiceConfig{}, err
	}
	svc, err := NewService(w, serviceName)
	if err != nil {
		return ServiceConfig{}, err
	}
	return svc.ReadConfig()
}

func SaveServiceConfig(workspaceDir string, serviceName string, conf ServiceConfig) error {
	svcConfigPath, err := getServiceConfigPathByLang(conf.Language)
	if err != nil {
		return err
	}
	path := filepath.Join(workspaceDir, svcConfigPath, serviceName)
	data, err := yaml.Marshal(&conf)
	if err != nil {
		return fmt.Errorf("failed to create service config: %w", err)
	}

	err = ioutil.WriteFile(filepath.Join(path, ServiceConfigName), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create service config: %w", err)
	}
	return nil
}

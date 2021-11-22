package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/service/lang"
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
	Language lang.ServiceLanguage `yaml:"-"`

	ServiceName string   `yaml:"service_name"`
	Maintainers []string `yaml:"maintainers"`

	OpenAPI ServiceOpenAPIConfig `yaml:"openapi,omitempty"`
}

func getServiceConfigPathByLang(language lang.ServiceLanguage) (string, error) {
	switch(language) {
	case lang.ServiceLanguageGo:
		return "go_services/cmd", nil
	case lang.ServiceLanguageJs:
		return "js_services", nil
	}
	return "", fmt.Errorf("no such language: %s", language)
}

func locateServiceConfigPath(workspaceDir string, serviceName string) (string, lang.ServiceLanguage, error) {
	for _, l := range lang.LanguagesList {
		subDir, _ := getServiceConfigPathByLang(l)
		configFile := filepath.Join(workspaceDir, subDir, serviceName, serviceConfigName)
		if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
			continue
		}
		return configFile, l, nil
	}
	return "", lang.ServiceLanguageUnknown, fmt.Errorf("service %s not found", serviceName)
}

func ReadServiceConfig(workspaceDir string, serviceName string) (ServiceConfig, error) {
	configFile, language, err := locateServiceConfigPath(workspaceDir, serviceName)
	if err != nil {
		return ServiceConfig{}, err
	}

	rawData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return ServiceConfig{}, fmt.Errorf("failed to read service config: %w", err)
	}

	var data ServiceConfig
	err = yaml.Unmarshal(rawData, &data)
	if err != nil {
		return ServiceConfig{}, fmt.Errorf("failed to read service config: %w", err)
	}
	data.Language = language

	return data, nil
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

	err = ioutil.WriteFile(filepath.Join(path, serviceConfigName), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create service config: %w", err)
	}
	return nil
}

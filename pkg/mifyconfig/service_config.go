package mifyconfig

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	ServiceConfigName = "service.mify.yaml"
)

type ServiceLanguage string

const (
	ServiceLanguageUnknown ServiceLanguage = "unknown"
	ServiceLanguageGo      ServiceLanguage = "go"
	ServiceLanguagePython  ServiceLanguage = "python"
	ServiceLanguageJs      ServiceLanguage = "js"
)

var LanguagesList = []ServiceLanguage{
	ServiceLanguageGo,
	ServiceLanguagePython,
	ServiceLanguageJs,
}

var (
	ErrNoSuchService = errors.New("no such service")
)

type ServiceOpenAPIClientConfig struct{}

type ServiceOpenAPIConfig struct {
	Clients map[string]ServiceOpenAPIClientConfig `yaml:"clients,omitempty"`
}

func (s ServiceOpenAPIConfig) HasClient(target string) bool {
	_, ok := s.Clients[target]
	return ok
}

type PostgresConfig struct {
	Enabled      bool   `yaml:"enabled"`
	DatabaseName string `yaml:"database_name,omitempty"`
}

type ComponentConfig struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type ComponentsConfig struct {
	Layout *ComponentConfig `yaml:"layout,omitempty"`
}

type ServiceConfig struct {
	Language ServiceLanguage `yaml:"language"`
	Template string          `yaml:"template,omitempty"`

	ServiceName string   `yaml:"service_name"`
	Maintainers []string `yaml:"maintainers"`
	Components ComponentsConfig `yaml:"components,omitempty"`

	OpenAPI  ServiceOpenAPIConfig `yaml:"openapi,omitempty"`
	Postgres PostgresConfig       `yaml:"postgres,omitempty"`

	IsExternal bool `yaml:"-"`
}

func ReadServiceCfg(path string) (*ServiceConfig, error) {
	wrapErr := func(err error) error {
		return fmt.Errorf("failed to read service config: %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, wrapErr(ErrNoSuchService)
	}

	rawData, err := os.ReadFile(path)
	if err != nil {
		return nil, wrapErr(err)
	}

	var data ServiceConfig
	err = yaml.Unmarshal(rawData, &data)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &data, nil
}

func MakeDefaultComponent() *ComponentConfig {
	return &ComponentConfig{
		Enabled: true,
	}
}

func tryReadExternalService(workspaceDir, serviceName string) (*ServiceConfig, error) {
	wrapErr := func(err error) error {
		return fmt.Errorf("failed to read external service config: %w", err)
	}
	externalSchemaPath := path.Join(workspaceDir, "schemas", "mify-external", serviceName, "api", "api.yaml")
	if _, err := os.Stat(externalSchemaPath); os.IsNotExist(err) {
		return nil, wrapErr(ErrNoSuchService)
	}
	return &ServiceConfig{
		ServiceName: serviceName,
		IsExternal:  true,
	}, nil
}

// TODO: remove (use ReadServiceCfg)
func ReadServiceConfig(workspaceDir string, serviceName string) (ServiceConfig, error) {
	schemaDir := path.Join(workspaceDir, "schemas", serviceName)
	path := filepath.Join(schemaDir, ServiceConfigName)

	cfg, err := ReadServiceCfg(path)
	if err == nil {
		return *cfg, nil
	}
	if err != nil && !errors.Is(err, ErrNoSuchService) {
		return ServiceConfig{}, err
	}
	cfg, err = tryReadExternalService(workspaceDir, serviceName)
	if err != nil {
		return ServiceConfig{}, err
	}

	return *cfg, nil
}

// Legacy, try to use Dump instead
func SaveServiceConfig(workspaceDir string, serviceName string, conf ServiceConfig) error {
	schemaDir := path.Join(workspaceDir, "schemas", serviceName)
	path := filepath.Join(schemaDir, ServiceConfigName)
	err := conf.Dump(path)
	if err != nil {
		return err
	}

	return nil
}

func (conf ServiceConfig) Dump(path string) error {
	wrapErr := func(err error) error {
		return fmt.Errorf("failed to dump service config: %w", err)
	}

	data, err := yaml.Marshal(&conf)
	if err != nil {
		return wrapErr(err)
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return wrapErr(err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return wrapErr(err)
	}
	return nil
}

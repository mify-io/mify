package mifyconfig

import (
	"fmt"
	"io/ioutil"
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
	ServiceLanguageJs      ServiceLanguage = "js"
)

var LanguagesList = []ServiceLanguage{
	ServiceLanguageGo,
	ServiceLanguageJs,
}

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

type ServiceConfig struct {
	Language ServiceLanguage `yaml:"language"`

	ServiceName string   `yaml:"service_name"`
	Maintainers []string `yaml:"maintainers"`

	OpenAPI  ServiceOpenAPIConfig `yaml:"openapi,omitempty"`
	Postgres PostgresConfig       `yaml:"postgres,omitempty"`
}

func ReadServiceCfg(path string) (*ServiceConfig, error) {
	wrapErr := func(err error) error {
		return fmt.Errorf("failed to read service config: %w", err)
	}

	rawData, err := ioutil.ReadFile(path)
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

// TODO: remove (use ReadServiceCfg)
func ReadServiceConfig(workspaceDir string, serviceName string) (ServiceConfig, error) {
	schemaDir := path.Join(workspaceDir, "schemas", serviceName)
	path := filepath.Join(schemaDir, ServiceConfigName)

	cfg, err := ReadServiceCfg(path)
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

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return wrapErr(err)
	}
	return nil
}

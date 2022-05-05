package cloudconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type CloudConfig struct {
	PublicHostname string `yaml:"public_hostname"`
	Path           string `yaml:"path"`
}

type ServiceCloudConfigDomain struct {
	CustomHostname string `yaml:"custom_hostname,omitempty"`
	Path           string `yaml:"path,omitempty"`
}

type EnvVariable struct {
	Value            *string           `yaml:"value,omitempty"`
	SecretName       string            `yaml:"secret_name,omitempty"`
	ValuePerEnv      map[string]string `yaml:"value_per_env,omitempty"`
	SecretNamePerEnv map[string]string `yaml:"secret_name_per_env,omitempty"`
}

type ServiceCloudConfig struct {
	Domain  ServiceCloudConfigDomain `yaml:"domain,omitempty"`
	Presets map[string]string        `yaml:"presets,omitempty"`
	EnvVars map[string]EnvVariable   `yaml:"env_vars,omitempty"`
	Publish bool                     `yaml:"publish,omitempty"`
}

func ReadServiceCloudCfg(path string) (*ServiceCloudConfig, error) {
	wrapErr := func(err error) error {
		return fmt.Errorf("failed to read cloud config: %w", err)
	}

	rawData, err := os.ReadFile(path)
	if err != nil {
		return nil, wrapErr(err)
	}

	var data ServiceCloudConfig
	err = yaml.Unmarshal(rawData, &data)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &data, nil
}

func (conf *ServiceCloudConfig) WriteToFile(path string) error {
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

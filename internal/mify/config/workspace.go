package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	workspaceConfigName = "workspace.mify.yaml"
)

type WorkspaceConfig struct {
	WorkspaceName string `yaml:"workspace_name"`
}

func ReadWorkspaceConfig() (WorkspaceConfig, error) {
	workspaceConfFile, err := ioutil.ReadFile(workspaceConfigName)

	if errors.Is(err, os.ErrNotExist) {
		return WorkspaceConfig{}, fmt.Errorf("workspace config not found, probably current directory is not a workspace")
	}
	if err != nil {
		return WorkspaceConfig{}, err
	}

	var data WorkspaceConfig

	err = yaml.Unmarshal(workspaceConfFile, &data)
	if err != nil {
		return WorkspaceConfig{}, fmt.Errorf("failed to read workspace config: %w", err)
	}

	return data, nil
}

func SaveWorkspaceConfig(path string, conf WorkspaceConfig) error {
	data, err := yaml.Marshal(&conf)
	if err != nil {
		return fmt.Errorf("failed to create workspace config: %w", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", path, workspaceConfigName), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create workspace config: %w", err)
	}
	return nil
}

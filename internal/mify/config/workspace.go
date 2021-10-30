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
	workspaceConfigName = "workspace.mify.yaml"
)

type WorkspaceConfig struct {
	WorkspaceName string `yaml:"workspace_name"`
	GitHost       string `yaml:"git_host"`
	GitNamespace  string `yaml:"git_namespace"`
	GitRepository string `yaml:"git_repository"`
}

func ReadWorkspaceConfig(path string) (WorkspaceConfig, error) {
	workspaceConfFile, err := ioutil.ReadFile(filepath.Join(path, workspaceConfigName))

	if errors.Is(err, os.ErrNotExist) {
		return WorkspaceConfig{}, fmt.Errorf("workspace config not found at path: %s", path)
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

func FindWorkspaceConfigPath() (string, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for curDir != "/" {
		if _, err := os.Stat(filepath.Join(curDir, workspaceConfigName)); err == nil {
			return curDir, nil
		}
		curDir = filepath.Dir(curDir)
	}
	return "", fmt.Errorf("unable to find workspace.mify.yaml, current or any parent directory is not a workspace")
}

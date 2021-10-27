package workspace

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
}

func CreateWorkspace(dir string, name string) error {
	fmt.Printf("creating workspace %s\n", name)

	context := Context{
		Name:     name,
		BasePath: filepath.Join(dir, name),
		GoRoot:   filepath.Join(dir, "go_services"),
	}

	if err := RenderTemplateTree(context); err != nil {
		return err
	}

	if err := createYaml(filepath.Join(dir, name)); err != nil {
		return err
	}
	return nil
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

// private

func createYaml(dir string) error {
	fmt.Printf("creating yaml in %s\n", dir)

	conf := WorkspaceConfig{
		WorkspaceName: dir,
	}

	data, err := yaml.Marshal(&conf)
	if err != nil {
		return fmt.Errorf("failed to create workspace config: %w", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, workspaceConfigName), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create workspace config: %w", err)
	}

	return nil
}

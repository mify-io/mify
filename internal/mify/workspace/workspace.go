package workspace

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	goModTemplate = "module %s%s"
)

const (
	workspaceConfigName = "workspace.mify.yaml"
)

type WorkspaceConfig struct {
	WorkspaceName string `yaml:"workspace_name"`
}

func CreateWorkspace(name string) error {
	fmt.Printf("creating workspace %s\n", name)

	if err := createHier(name); err != nil {
		return err
	}
	if err := createYaml(name); err != nil {
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

func createHier(dir string) error {
	fmt.Printf("creating hierarchy in %s\n", dir)

	err := os.Mkdir(dir, 0755)
	if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("failed to create base directory: please remove file or directory with the same name")
	}
	if err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}

	// TODO: README.md
	// TODO: init git
	// TODO: vcs specific files (.gitowners, .gitignore)
	basePaths := []string{
		"schemas",
		"frontend",
		"backend/cmd",
		"backend/internal/pkg",
		"backend/pkg",
	}
	for _, path := range basePaths {
		err = os.MkdirAll(fmt.Sprintf("%s/%s", dir, path), 0755)
		if err != nil {
			return fmt.Errorf("failed to create %s directory: %w", path, err)
		}
	}

	goModRendered := fmt.Sprintf(goModTemplate, "repo.com/namespace/", dir)
	err = ioutil.WriteFile(fmt.Sprintf("%s/backend/go.mod", dir), []byte(goModRendered), 0644)
	if err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}

	return nil
}

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

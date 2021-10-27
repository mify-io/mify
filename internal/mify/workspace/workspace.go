package workspace

import (
	"fmt"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/config"
)

const (
	workspaceConfigName = "workspace.mify.yaml"
)

type WorkspaceConfig struct {
	WorkspaceName string `yaml:"workspace_name"`
}

func CreateWorkspace(dir string, name string) error {
// var (
	// goModTemplate = "module %s%s"
// )

// func CreateWorkspace(name string) error {
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

// private

func createYaml(dir string) error {
	fmt.Printf("creating yaml in %s\n", dir)

	conf := config.WorkspaceConfig{
		WorkspaceName: dir,
	}
	return config.SaveWorkspaceConfig(dir, conf)
}

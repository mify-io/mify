package workspace

import (
	"fmt"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/core"
)

func CreateWorkspace(ctx *core.Context, dir string, name string) error {
	fmt.Printf("Creating workspace: %s\n", name)

	context := Context{
		Name:     name,
		BasePath: filepath.Join(dir, name),
		GoRoot:   filepath.Join(dir, "go_services"),
	}

	if err := RenderTemplateTree(ctx, context); err != nil {
		return err
	}

	if err := createYaml(name, filepath.Join(dir, name)); err != nil {
		return err
	}
	return nil
}

// private

func createYaml(name, dir string) error {
	conf := config.WorkspaceConfig{
		WorkspaceName: name,
		GitHost: "example.com",
		GitNamespace: "namespace",
	}
	return config.SaveWorkspaceConfig(dir, conf)
}


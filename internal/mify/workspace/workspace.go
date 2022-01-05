package workspace

import (
	"fmt"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/pkg/mifyconfig"
)

func CreateWorkspace(ctx *core.Context, dir string, name string) error {
	fmt.Printf("Creating workspace: %s\n", name)

	context := Context{
		Name:     name,
		BasePath: filepath.Join(dir, name),
		GoRoot:   filepath.Join(dir, mifyconfig.GoServicesRoot),
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
	conf := mifyconfig.WorkspaceConfig{
		WorkspaceName: name,
		GitHost: "example.com",
		GitNamespace: "namespace",
	}
	return mifyconfig.SaveWorkspaceConfig(dir, conf)
}


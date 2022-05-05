package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace/mutators"
)

const (
	SchemasDirName    = "schemas"
)

func CreateWorkspace(mutContext *mutators.MutatorContext, dirAbsPath string, name string) error {
	wrapErr := func(err error) error {
		return fmt.Errorf("error while creating workspace: %w", err)
	}

	fmt.Printf("Creating workspace: %s\n", name)

	baseAbsPath := filepath.Join(dirAbsPath, name)
	if err := os.MkdirAll(baseAbsPath, os.ModePerm); err != nil {
		return wrapErr(err)
	}

	schemasAbsPath := filepath.Join(baseAbsPath, SchemasDirName)
	if err := os.MkdirAll(schemasAbsPath, os.ModePerm); err != nil {
		return wrapErr(err)
	}

	gitignoreAbsPath := filepath.Join(baseAbsPath, ".gitignore")
	if _, err := os.Create(gitignoreAbsPath); err != nil {
		return wrapErr(err)
	}

	if err := createYaml(name, baseAbsPath); err != nil {
		return wrapErr(err)
	}
	return nil
}

func createYaml(name, dir string) error {
	conf := mifyconfig.WorkspaceConfig{
		WorkspaceName: name,
		GitHost:       "example.com",
		GitNamespace:  "namespace",
	}
	return mifyconfig.SaveWorkspaceConfig(dir, conf)
}

package workspace

import (
	"fmt"
	"path/filepath"
)

const goServices = "go_services"

type Context struct {
	Name     string
	BasePath string
	GoRoot   string // Path to go_services
}

func InitContext(workspacePath string) (Context, error) {
	fmt.Printf("workspacePath %s\n", workspacePath)
	fmt.Printf("go root %s\n", filepath.Join(workspacePath, goServices))

	res := Context{
		Name:     filepath.Base(workspacePath), // TODO: validate
		BasePath: workspacePath,
		GoRoot:   filepath.Join(workspacePath, goServices),
	}

	return res, nil
}

package service

import (
	"fmt"

	// "github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/workspace"
)

func CreateService(wspConext workspace.Context, name string) error {
	fmt.Printf("creating service %s\n", name)

	context := Context{
		ServiceName: name,
		Workspace:   wspConext,

	}
	// _, err := config.ReadWorkspaceConfig()
	// if err != nil {
		// return err
	// }

	if err := RenderTemplateTree(context); err != nil {
		return err
	}

	// _, err := workspace.ReadWorkspaceConfig()
	// if err != nil {
	// 	return err
	// }

	// if err := createServiceYaml(name); err != nil {
	// 	return err
	// }

	return nil
}

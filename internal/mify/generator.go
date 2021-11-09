package mify

import (
	"context"

	"github.com/chebykinn/mify/internal/mify/service"
	"github.com/chebykinn/mify/internal/mify/util/docker"
	"github.com/chebykinn/mify/internal/mify/workspace"
)

func CreateWorkspace(basePath string, name string) error {
	return workspace.CreateWorkspace(basePath, name)
}

func CreateService(workspacePath string, name string) error {
	workspaceContext, err := workspace.InitContext(workspacePath)
	if err != nil {
		return err
	}

	return service.CreateService(workspaceContext, name)
}

func ServiceGenerate(workspacePath string, name string) error {
	workspaceContext, err := workspace.InitContext(workspacePath)
	if err != nil {
		return err
	}

	return service.Generate(workspaceContext, name)
}

func Cleanup() error {
	err := docker.Cleanup(context.Background())
	if err != nil {
		return err
	}

	return nil
}

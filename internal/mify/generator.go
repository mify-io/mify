package mify

import (
	"github.com/chebykinn/mify/internal/mify/service"
	"github.com/chebykinn/mify/internal/mify/workspace"
)

func CreateWorkspace(name string) error {
	return workspace.CreateWorkspace(name)
}

func CreateService(name string) error {
	return service.CreateService(name)
}

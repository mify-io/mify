package service

import (
	"github.com/chebykinn/mify/internal/mify/workspace"
)

type Context struct {
	ServiceName string
	Workspace   workspace.Context
}

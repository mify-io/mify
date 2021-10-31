package service

import (
	"github.com/chebykinn/mify/internal/mify/workspace"
)

type Context struct {
	ServiceName string
	Repository  string
	GoModule    string
	Workspace   workspace.Context
}

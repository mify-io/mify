package service

import (
	"github.com/chebykinn/mify/internal/mify/workspace"
	"github.com/chebykinn/mify/pkg/mifyconfig"
)

type OpenAPIClientContext struct {
	ClientName string
	PackageName string
	PrivateFieldName string
	PublicMethodName string
}

type OpenAPIContext struct {
	Clients []OpenAPIClientContext
}

type Context struct {
	ServiceName string
	Repository  string
	Language    mifyconfig.ServiceLanguage
	GoModule    string
	Workspace   workspace.Context
	OpenAPI     OpenAPIContext
	ServiceList []string
}

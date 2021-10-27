package service

import (
	"strings"

	"github.com/chebykinn/mify/internal/mify/core"
)

const (
	serviceNamePlaceholder = "_svc_"
)

func transformPath(context interface{}, path string) (string, error) {
	tContext := context.(Context)

	path = strings.ReplaceAll(path, serviceNamePlaceholder, tContext.ServiceName)

	return path, nil
}

func RenderTemplateTree(context Context) error {
	params := core.RenderParams{
		TemplatesPath:   "assets/tpl/go_service",
		TargetPath:      context.Workspace.GoRoot,
		PathTransformer: transformPath,
	}
	return core.RenderTemplateTree(context, params)
}

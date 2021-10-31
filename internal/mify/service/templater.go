package service

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/chebykinn/mify/internal/mify/core"
)

const (
	serviceNamePlaceholder = "#svc#"
)

func transformPath(context interface{}, path string) (string, error) {
	tContext := context.(Context)

	path = strings.ReplaceAll(path, serviceNamePlaceholder, tContext.ServiceName)

	return path, nil
}

func RenderTemplateTree(context Context) error {
	funcMap := template.FuncMap{
		"svcUserCtxName": func(context Context) string {
			return fmt.Sprintf("%s%s", strings.Title(context.ServiceName), "Context")
		},
	}
	params := core.RenderParams{
		TemplatesPath:   "tpl/go_service",
		TargetPath:      context.Workspace.BasePath,
		PathTransformer: transformPath,
		FuncMap:         funcMap,
	}
	return core.RenderTemplateTree(context, params)
}

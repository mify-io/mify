package service

import (
	"fmt"
	"path/filepath"
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

func RenderTemplateTree(ctx *core.Context, context Context) error {
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
	return core.RenderTemplateTree(ctx, context, params)
}

func RenderTemplateTreeSubPath(ctx *core.Context, context Context, templateSubPath string) error {
	funcMap := template.FuncMap{
		"svcUserCtxName": func(context Context) string {
			return fmt.Sprintf("%s%s", strings.Title(context.ServiceName), "Context")
		},
	}
	targetSubPath, err := transformPath(context, templateSubPath)
	if err != nil {
		return err
	}
	params := core.RenderParams{
		TemplatesPath:   filepath.Join("tpl/go_service", templateSubPath),
		TargetPath:      filepath.Join(context.Workspace.BasePath, targetSubPath),
		PathTransformer: transformPath,
		FuncMap:         funcMap,
	}
	return core.RenderTemplateTree(ctx, context, params)
}

package service

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/pkg/mifyconfig"
)

const (
	serviceNamePlaceholder = "#svc#"
)

func transformPath(context interface{}, path string) (string, error) {
	tContext := context.(Context)

	path = strings.ReplaceAll(path, serviceNamePlaceholder, tContext.ServiceName)

	return path, nil
}

func getLanguageTemplatePath(context Context) (string, error) {
	switch context.Language {
	case mifyconfig.ServiceLanguageGo:
		return "tpl/go_service", nil
	case mifyconfig.ServiceLanguageJs:
		return "tpl/js_service", nil
	}
	return "", fmt.Errorf("no such language: %s", context.Language)
}

func RenderTemplateTree(ctx *core.Context, context Context) error {
	funcMap := template.FuncMap{
		"svcUserCtxName": func(context Context) string {
			return fmt.Sprintf("%s%s", strings.Title(context.ServiceName), "Context")
		},
	}
	templatesPath, err := getLanguageTemplatePath(context)
	if err != nil {
		return err
	}

	params := core.RenderParams{
		TemplatesPath:   templatesPath,
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
	templatesPath, err := getLanguageTemplatePath(context)
	if err != nil {
		return err
	}

	targetSubPath, err := transformPath(context, templateSubPath)
	if err != nil {
		return err
	}
	params := core.RenderParams{
		TemplatesPath:   filepath.Join(templatesPath, templateSubPath),
		TargetPath:      filepath.Join(context.Workspace.BasePath, targetSubPath),
		PathTransformer: transformPath,
		FuncMap:         funcMap,
	}
	return core.RenderTemplateTree(ctx, context, params)
}

package layout

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

const (
	serviceNamePlaceholder = "#svc#"
)

func execute(ctx *gencontext.GenContext) error {
	if err := renderServiceTemplateTree(ctx, tpl.NewServiceModel(ctx)); err != nil {
		return fmt.Errorf("error while rendering service: %w", err)
	}

	return nil
}

func renderServiceTemplateTree(ctx *gencontext.GenContext, model *tpl.ServiceModel) error {
	funcMap := template.FuncMap{
		"svcUserCtxName": func(model tpl.ServiceModel) string {
			return fmt.Sprintf("%s%s", strings.Title(ctx.GetServiceName()), "Context")
		},
	}
	templatesPath, err := getLanguageTemplatePath(ctx)
	if err != nil {
		return err
	}

	params := RenderParams{
		TemplatesPath:   templatesPath,
		TargetPath:      ctx.GetWorkspace().BasePath,
		PathTransformer: serviceTransformPath,
		FuncMap:         funcMap,
	}
	return RenderTemplateTree(ctx, model, params)
}

func serviceTransformPath(model interface{}, path string) (string, error) {
	tmodel := model.(*tpl.ServiceModel)

	path = strings.ReplaceAll(path, serviceNamePlaceholder, tmodel.ServiceName)

	return path, nil
}

func getLanguageTemplatePath(ctx *gencontext.GenContext) (string, error) {
	mifySchema := ctx.MustGetMifySchema()

	switch mifySchema.Language {
	case mifyconfig.ServiceLanguageGo:
		return "assets/go_service", nil
	case mifyconfig.ServiceLanguageJs:
		return "assets/js_service", nil
	}
	if len(mifySchema.Language) == 0 {
		return "", fmt.Errorf("missing language in service.mify.yaml")
	}
	return "", fmt.Errorf("no such language: %s", mifySchema.Language)
}
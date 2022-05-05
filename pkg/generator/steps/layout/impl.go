package layout

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl"
	tplnew "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	serviceNamePlaceholder = "#svc#"
)

func execute(ctx *gencontext.GenContext) error {
	if ctx.GetMifySchema() != nil && ctx.MustGetMifySchema().Language == mifyconfig.ServiceLanguageGo {
		if err := renderServiceTemplateTree(ctx, tpl.NewServiceModel(ctx)); err != nil {
			return fmt.Errorf("error while rendering service: %w", err)
		}
	}

	if err := renderNew(ctx); err != nil {
		return err
	}

	return nil
}

func renderServiceTemplateTree(ctx *gencontext.GenContext, model *tpl.ServiceModel) error {
	funcMap := template.FuncMap{
		"svcUserCtxName": func(model tpl.ServiceModel) string {
			caser := cases.Title(language.AmericanEnglish)
			return fmt.Sprintf("%s%s", caser.String(ctx.GetServiceName()), "Context")
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
		return "assets/go-service", nil
	case mifyconfig.ServiceLanguageJs:
		return "assets/js-service", nil
	}
	if len(mifySchema.Language) == 0 {
		return "", fmt.Errorf("missing language in service.mify.yaml")
	}
	return "", fmt.Errorf("no such language: %s", mifySchema.Language)
}

func renderNew(ctx *gencontext.GenContext) error {

	mifySchema := ctx.GetMifySchema()
	if mifySchema == nil {
		return nil
	}

	switch mifySchema.Language {
	case mifyconfig.ServiceLanguageGo:
		if err := tplnew.RenderGo(ctx); err != nil {
			return fmt.Errorf("can't render go files: %w", err)
		}
	case mifyconfig.ServiceLanguageJs:
		if err := tplnew.RenderJs(ctx); err != nil {
			return fmt.Errorf("can't render js files: %w", err)
		}
	case mifyconfig.ServiceLanguagePython:
		if err := tplnew.RenderPy(ctx); err != nil {
			return fmt.Errorf("can't render python files: %w", err)
		}
	}

	return nil
}

package layout

import (
	_ "embed"
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

func execute(ctx *gencontext.GenContext) error {
	if err := render(ctx); err != nil {
		return err
	}

	return nil
}

func render(ctx *gencontext.GenContext) error {
	mifySchema := ctx.GetMifySchema()
	if mifySchema == nil {
		return nil
	}

	switch mifySchema.Language {
	case mifyconfig.ServiceLanguageGo:
		if err := tpl.RenderGo(ctx); err != nil {
			return fmt.Errorf("can't render go files: %w", err)
		}
	case mifyconfig.ServiceLanguageJs:
		if err := tpl.RenderJs(ctx); err != nil {
			return fmt.Errorf("can't render js files: %w", err)
		}
	case mifyconfig.ServiceLanguagePython:
		if err := tpl.RenderPy(ctx); err != nil {
			return fmt.Errorf("can't render python files: %w", err)
		}
	}

	return nil
}

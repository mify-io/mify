package layout

import (
	_ "embed"
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	tplnew "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

func execute(ctx *gencontext.GenContext) error {
	if err := renderNew(ctx); err != nil {
		return err
	}

	return nil
}

func renderNew(ctx *gencontext.GenContext) error {

	mifySchema := ctx.GetMifySchema()
	if mifySchema == nil {
		return nil
	}
	if err := tplnew.RenderWorkspace(ctx); err != nil {
		return fmt.Errorf("can't render workspace files: %w", err)
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

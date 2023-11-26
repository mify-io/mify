package layout

import (
	_ "embed"

	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type LayoutStep struct {
}

func NewLayoutStep() LayoutStep {
	return LayoutStep{}
}

func (s LayoutStep) Name() string {
	return "layout"
}

func (s LayoutStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	if ctx.GetMifySchema() != nil && !ctx.MustGetMifySchema().Components.Layout.Enabled {
		ctx.Logger.Info("skipping disabled step")
		return core.Done, nil
	}
	if err := execute(ctx); err != nil {
		return core.Done, err
	}

	return core.Done, nil
}

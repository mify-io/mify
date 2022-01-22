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
	if err := execute(ctx); err != nil {
		return core.Done, err
	}

	return core.Done, nil
}

package layout

import (
	_ "embed"

	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type MifyGeneratedStep struct {
}

func NewMifyGeneratedStep() MifyGeneratedStep {
	return MifyGeneratedStep{}
}

func (s MifyGeneratedStep) Name() string {
	return "mify-generated"
}

func (s MifyGeneratedStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	if err := execute(ctx); err != nil {
		return core.Done, err
	}

	return core.Done, nil
}

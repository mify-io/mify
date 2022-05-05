package prepare

import (
	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type PrepareStep struct {
}

func NewPrepareStep() PrepareStep {
	return PrepareStep{}
}

func (s PrepareStep) Name() string {
	return "prepare"
}

func (s PrepareStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	if err := execute(ctx); err != nil {
		return core.Done, err
	}

	return core.Done, nil
}

package devrunner

import (
	_ "embed"

	"github.com/chebykinn/mify/pkg/generator/core"
	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
)

type DevRunnerStep struct {
}

func NewDevRunnerStep() DevRunnerStep {
	return DevRunnerStep{}
}

func (s DevRunnerStep) Name() string {
	return "DevRunner"
}

func (s DevRunnerStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	if err := execute(ctx); err != nil {
		return core.Done, err
	}

	return core.Done, nil
}

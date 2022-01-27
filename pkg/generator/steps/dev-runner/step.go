package devrunner

import (
	_ "embed"

	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/workspace"
)

type DevRunnerStep struct {
}

func NewDevRunnerStep() DevRunnerStep {
	return DevRunnerStep{}
}

func (s DevRunnerStep) Name() string {
	return "dev-runner"
}

func (s DevRunnerStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	if ctx.GetServiceName() != workspace.DevRunnerName {
		return core.Done, nil
	}

	if err := execute(ctx); err != nil {
		return core.Done, err
	}

	return core.Done, nil
}

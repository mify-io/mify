package core

import (
	"context"

	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
	"github.com/chebykinn/mify/pkg/workspace"
)

type Pipeline struct {
	steps []Step
}

func (p Pipeline) Execute(
	goContext context.Context,
	serviceName string,
	workspaceDescription workspace.Description) error {

	genContext := gencontext.NewGenContext(goContext, serviceName, workspaceDescription)
	for _, step := range p.steps {
		if err := step.Execute(genContext); err != nil {
			return err
		}
	}

	return nil
}

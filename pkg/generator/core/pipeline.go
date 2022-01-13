package core

import (
	"context"
	"fmt"

	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
	"github.com/chebykinn/mify/pkg/workspace"
)

const (
	maxRepeatsCount = 20
)

type Pipeline struct {
	steps []Step
}

func (p Pipeline) Execute(
	goContext context.Context,
	serviceName string,
	workspaceDescription workspace.Description) error {

	shouldRepeat := true
	iteration := 0
	for shouldRepeat {
		iteration++
		if iteration == maxRepeatsCount {
			return fmt.Errorf("max number %d of pipeline execution repeats has been reached", maxRepeatsCount)
		}

		shouldRepeat = false
		genContext := gencontext.NewGenContext(goContext, serviceName, workspaceDescription)
		for _, step := range p.steps {
			result, err := step.Execute(genContext)
			if err != nil {
				return err
			}

			if result == RepeatAll {
				shouldRepeat = true
				break
			}
		}
	}

	return nil
}

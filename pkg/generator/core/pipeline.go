package core

import (
	"context"
	"fmt"

	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
	"github.com/chebykinn/mify/pkg/mifyconfig"
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

		serviceCfg, err := mifyconfig.ReadServiceConfig(workspaceDescription.BasePath, serviceName)
		if err != nil {
			return err
		}
		genContext := gencontext.NewGenContext(goContext, serviceName, workspaceDescription, serviceCfg)

		for _, step := range p.steps {
			genContext.Logger.Println(fmt.Sprintf("Starting step '%s'", step.Name()))
			result, err := step.Execute(genContext)
			if err != nil {
				return fmt.Errorf("Step '%s' failed with error: '%w'", step.Name(), err)
			}

			if result == RepeatAll {
				shouldRepeat = true
				break
			}
		}
	}

	return nil
}

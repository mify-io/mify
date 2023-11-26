package core

import (
	"context"
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/workspace"
	"go.uber.org/zap"
)

const (
	maxRepeatsCount = 20
)

type StepExecResult struct {
	SeqNo int
	StepName string
	Error error
}

type Pipeline struct {
	steps []Step
}

// Each step can be returned as completed several times, because of possible pipeline rerun
// Step with seq_no = -1 means some error inside pipeline, not step
func (p Pipeline) Execute(
	goContext context.Context,
	serviceName string,
	workspaceDescription workspace.Description,
	migrate bool,
	forceRegeneration bool,
	mifyVersion string,
	verboseOutput bool,
	outChan chan StepExecResult) {

	shouldRepeat := true
	iteration := 0
	for shouldRepeat {
		iteration++
		if iteration == maxRepeatsCount {
			outChan <- StepExecResult{
				SeqNo: -1,
				StepName: "",
				Error: fmt.Errorf("max number %d of pipeline execution repeats has been reached", maxRepeatsCount),
			}
		}

		shouldRepeat = false

		genContext, err := gencontext.NewGenContext(
			goContext, serviceName, workspaceDescription,
			migrate, forceRegeneration, mifyVersion, verboseOutput)
		if err != nil {
			outChan <- StepExecResult{
				SeqNo: -1,
				StepName: "",
				Error: fmt.Errorf("can't initialize generate pipeline: %w", err),
			}
			break
		}

		logger := genContext.Logger
		for stepSeqNo, step := range p.steps {
			genContext.Logger = logger.With(
				zap.String("step_name", step.Name()),
			)
			execRes := StepExecResult{
				SeqNo: stepSeqNo,
				StepName:  step.Name(),
			}
			outChan <- execRes

			genContext.Logger.Infof("Starting step")
			result, err := step.Execute(genContext)

			if err != nil {
				execRes.Error = fmt.Errorf("step '%s' failed with error: %w", step.Name(), err)
			}
			genContext.Logger.Infof("Finished step")

			outChan <- execRes

			if result == RepeatAll {
				shouldRepeat = true
				break
			}
		}
	}
}

func (p Pipeline) Size() int {
	return len(p.steps)
}

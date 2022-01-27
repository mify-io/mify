package core

import (
	"context"
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/workspace"
)

const (
	maxRepeatsCount = 20
)

type StepExecResult struct {
	SeqNo int
	Step  *Step
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
	outChan chan StepExecResult) {

	shouldRepeat := true
	iteration := 0
	for shouldRepeat {
		iteration++
		if iteration == maxRepeatsCount {
			outChan <- StepExecResult{
				SeqNo: -1,
				Step:  nil,
				Error: fmt.Errorf("max number %d of pipeline execution repeats has been reached", maxRepeatsCount),
			}
		}

		shouldRepeat = false

		genContext := gencontext.NewGenContext(goContext, serviceName, workspaceDescription)

		for stepSeqNo, step := range p.steps {
			genContext.Logger.Infof("Starting step '%s'", step.Name())
			result, err := step.Execute(genContext)

			execRes := StepExecResult{
				SeqNo: stepSeqNo,
				Step:  &step,
			}

			if err != nil {
				execRes.Error = fmt.Errorf("step '%s' failed with error: %w", step.Name(), err)
			}

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

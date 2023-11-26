package status

import (
	"sync"

	"github.com/mify-io/mify/pkg/generator/core"
	"github.com/vbauerster/mpb/v7/decor"
)

type GeneratorCliProgress struct {
	progressBar       *ProgressBar
	lastCompletedStep core.StepExecResult
	mu sync.RWMutex
}

func (pg *GeneratorCliProgress) updateStatus(s decor.Statistics) string {
	pg.mu.RLock()
	stepName := pg.lastCompletedStep.StepName
	spinner := pg.progressBar.Spinner()
	pg.mu.RUnlock()
	return spinner + " running: [" + stepName + "] "
}

func NewGeneratorCliProgress(pipeline core.Pipeline) *GeneratorCliProgress {
	pg := GeneratorCliProgress{
		lastCompletedStep: core.StepExecResult{},
	}

	pg.progressBar = NewProgressBar(pg.updateStatus)
	pg.progressBar.Create(int64(pipeline.Size()))

	return &pg
}

func (pg *GeneratorCliProgress) ReportStep(stepExecResult core.StepExecResult) {
	pg.mu.Lock()
	if pg.lastCompletedStep.StepName == "" || pg.lastCompletedStep.SeqNo < stepExecResult.SeqNo {
		pg.lastCompletedStep = stepExecResult
		pg.progressBar.Increment()
	}
	pg.mu.Unlock()
}

func (pg *GeneratorCliProgress) Wait() {
	pg.progressBar.Wait()
}

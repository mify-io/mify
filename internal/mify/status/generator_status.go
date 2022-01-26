package status

import (
	"github.com/mify-io/mify/pkg/generator/core"
	"github.com/vbauerster/mpb/v7/decor"
)

type GeneratorCliProgress struct {
	progressBar       *ProgressBar
	lastCompletedStep *core.StepExecResult
}

func (pg *GeneratorCliProgress) updateStatus(s decor.Statistics) string {
	return pg.progressBar.Spinner() + " running: [" + (*pg.lastCompletedStep.Step).Name() + "] "
}

func NewGeneratorCliProgress(pipeline core.Pipeline) *GeneratorCliProgress {
	pg := GeneratorCliProgress{
		lastCompletedStep: nil,
	}

	pg.progressBar = NewProgressBar(pg.updateStatus)
	pg.progressBar.Create(int64(pipeline.Size()))

	return &pg
}

func (pg *GeneratorCliProgress) ReportStep(stepExecResult *core.StepExecResult) {
	if pg.lastCompletedStep == nil || pg.lastCompletedStep.SeqNo < stepExecResult.SeqNo {
		pg.lastCompletedStep = stepExecResult
		pg.progressBar.Increment()
	}
}

func (pg *GeneratorCliProgress) Wait() {
	pg.progressBar.Wait()
}

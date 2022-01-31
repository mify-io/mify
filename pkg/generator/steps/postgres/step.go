package postgres

import (
	_ "embed"

	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type PostgresStep struct {
}

func NewPostgresStep() PostgresStep {
	return PostgresStep{}
}

func (s PostgresStep) Name() string {
	return "postgres"
}

func (s PostgresStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	if err := execute(ctx); err != nil {
		return core.Done, err
	}

	return core.Done, nil
}

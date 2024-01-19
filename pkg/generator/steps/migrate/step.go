package migrate

import (
	_ "embed"

	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type MigrateStep struct {
}

func NewMigrateStep() MigrateStep {
	return MigrateStep{}
}

func (s MigrateStep) Name() string {
	return "migrate"
}

func (s MigrateStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	if !ctx.GetMigrate() {
		ctx.Logger.Info("called without --migrate, skipping")
		return core.Done, nil
	}
	if ctx.GetMifySchema() == nil {
		return core.Done, nil
	}
	if err := execute(ctx); err != nil {
		return core.Done, err
	}

	return core.Done, nil
}

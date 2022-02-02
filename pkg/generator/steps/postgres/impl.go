package postgres

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/postgres/tpl"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

func execute(ctx *gencontext.GenContext) error {
	if ctx.GetMifySchema() == nil {
		return nil
	}

	if !ctx.GetMifySchema().Postgres.Enabled {
		return nil
	}

	if ctx.GetMifySchema().Language == mifyconfig.ServiceLanguageGo {
		ctx.Logger.Infof("Will generated postgres for service")
		if err := tpl.RenderGo(ctx); err != nil {
			return err
		}
	}

	return nil
}

package app

import (
	_ "embed"
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed request_extra.go.tpl
var requestExtraTemplate string

//go:embed service_extra.go.tpl
var serviceExtraTemplate string

func Render(ctx *gencontext.GenContext) error {
	reqExtraModel := newRequestExtraModel(ctx)
	reqExtraPath := ctx.GetWorkspace().GetAppSubAbsPath(ctx.GetServiceName(), "request_extra.go")
	if err := render.RenderOrSkipTemplate(requestExtraTemplate, reqExtraModel, reqExtraPath); err != nil {
		return render.WrapError("request extra", err)
	}

	serviceExtraModel := newServiceExtraModel(ctx)
	serviceExtraPath := ctx.GetWorkspace().GetAppSubAbsPath(ctx.GetServiceName(), "service_extra.go")
	serviceExtraRelPath := ctx.GetWorkspace().GetAppSubRelPath(ctx.GetServiceName(), "service_extra.go")
	migrationSettings := render.MigrateSettings{
		Migrate:              ctx.GetMigrate(),
		HasUncommitedChanges: func() (bool, error) { return ctx.GetVcsIntegration().FileHasUncommitedChanges(serviceExtraRelPath) },
		Migrations:           []render.MigrationCallback{migrateContextToServiceExtra},
	}

	err := render.RenderOrMigrateTemplate(serviceExtraTemplate,
		serviceExtraModel, serviceExtraPath, migrationSettings)
	if err != nil {
		return render.WrapError(fmt.Sprintf("file %s", serviceExtraPath), err)
	}

	return nil
}

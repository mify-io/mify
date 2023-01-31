package core

import (
	_ "embed"
	"io/ioutil"
	"os"
	"path"
	"strings"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/migrate"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed request_context.go.tpl
var requestContextTemplate string

//go:embed service_context.go.tpl
var serviceContextTemplate string

//go:embed helpers.go.tpl
var helpersTemplate string

func maybeMigrate121(ctx *gencontext.GenContext, oldRequestContextPath string) {
	dat, err := ioutil.ReadFile(oldRequestContextPath)
	if os.IsNotExist(err) {
		return
	}

	if strings.Contains(string(dat), "mifyServiceContext *MifyServiceContext") {
		return // already migrated
	}

	migrate.MigrateSubstring(ctx, ctx.GetWorkspace().GetGoServicesAbsPath(), "ctx.MifyServiceContext", "", "ctx.ServiceContext()")
	migrate.MigrateSubstring(ctx, ctx.GetWorkspace().GetGoServicesAbsPath(), "ctx.ServiceExtra().(*app.ServiceExtra)", "package apputil", "apputil.GetServiceExtra(ctx.ServiceContext())/*TODO: you have to manually import package*/")
	migrate.MigrateSubstring(ctx, ctx.GetWorkspace().GetGoServicesAbsPath(), ".GetContext()", "", ".GoContext()")
	migrate.MigrateSubstring(ctx, ctx.GetWorkspace().GetGoServicesAbsPath(), ".RequestContext()", "", ".GoContext()")
}

func Render(ctx *gencontext.GenContext) error {
	requestContextModel := newRequestContextModel(ctx)
	requestContextPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetGoServiceGeneratedCoreRelPath(ctx.GetServiceName()), "request_context.go")
	if ctx.GetMigrate() {
		maybeMigrate121(ctx, requestContextPath)
	}
	if err := render.RenderTemplate(requestContextTemplate, requestContextModel, requestContextPath); err != nil {
		return render.WrapError("request_context", err)
	}

	serviceContextModel := newServiceContextModel(ctx)
	serviceContextPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetGoServiceGeneratedCoreRelPath(ctx.GetServiceName()), "service_context.go")
	if err := render.RenderTemplate(serviceContextTemplate, serviceContextModel, serviceContextPath); err != nil {
		return render.WrapError("service_context", err)
	}

	helpersModel := newHelpersModel(ctx)
	helpersPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetGoServiceGeneratedCoreRelPath(ctx.GetServiceName()), "helpers.go")
	if err := render.RenderTemplate(helpersTemplate, helpersModel, helpersPath); err != nil {
		return render.WrapError("helpers", err)
	}

	return nil
}

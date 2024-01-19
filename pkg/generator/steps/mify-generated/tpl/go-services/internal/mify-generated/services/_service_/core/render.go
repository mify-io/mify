package core

import (
	"embed"
	"io/ioutil"
	"os"
	"path"
	"strings"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/migrate"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed *.tpl
var templates embed.FS

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
	requestContextPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetGoServiceGeneratedCoreRelPath(ctx.GetServiceName()), "request_context.go")
	if ctx.GetMigrate() {
		maybeMigrate121(ctx, requestContextPath)
	}

	requestContextModel := render.NewModel(ctx, newRequestContextModel(ctx))
	serviceContextModel := render.NewModel(ctx, newServiceContextModel(ctx))
	helpersModel := render.NewModel(ctx, newHelpersModel(ctx))
	curPath := path.Join(
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetServicePath().Abs(),
		"core",
	)
	return render.RenderMany(templates,
		render.NewFile(ctx, path.Join(curPath, "request_context.go")).SetModel(requestContextModel),
		render.NewFile(ctx, path.Join(curPath, "service_context.go")).SetModel(serviceContextModel),
		render.NewFile(ctx, path.Join(curPath, "helpers.go")).SetModel(helpersModel),
	)
}

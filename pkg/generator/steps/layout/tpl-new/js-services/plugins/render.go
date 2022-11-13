package plugins

import (
	_ "embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed clients.js.tpl
var clientsTemplate string

func Render(ctx *gencontext.GenContext) error {
	curPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetJsServiceRelPath(ctx.GetServiceName()), "plugins")
	clientsModel := newClientsModel(ctx)
	clientsPath := path.Join(curPath, "clients.js")
	if err := render.RenderOrSkipTemplate(clientsTemplate, clientsModel, clientsPath); err != nil {
		return render.WrapError("clients.js", err)
	}

	return nil
}


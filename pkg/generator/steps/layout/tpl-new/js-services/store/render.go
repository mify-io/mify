package store

import (
	_ "embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed index.js.tpl
var indexTemplate string

//go:embed config.js.tpl
var configTemplate string

func Render(ctx *gencontext.GenContext) error {
	curPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetJsServiceRelPath(ctx.GetServiceName()), "store")
	indexModel := newIndexModel(ctx)
	indexPath := path.Join(curPath, "index.js")
	if err := render.RenderOrSkipTemplate(indexTemplate, indexModel, indexPath); err != nil {
		return render.WrapError("index.js", err)
	}

	configModel := newConfigModel(ctx)
	configPath := path.Join(curPath, "config.js")
	if err := render.RenderOrSkipTemplate(configTemplate, configModel, configPath); err != nil {
		return render.WrapError("config.js", err)
	}

	return nil
}

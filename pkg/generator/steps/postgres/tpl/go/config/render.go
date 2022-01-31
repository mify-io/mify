package config

import (
	_ "embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed postgres_config.go.tpl
var postgresConfigTemplate string

func Render(ctx *gencontext.GenContext) error {
	postgresConfigModel := NewPostgresConfigModel(ctx)
	// TODO: move path to description
	postgresConfigPath := path.Join(ctx.GetWorkspace().GetGoPostgresConfigAbsPath(), "config.go")
	if err := render.RenderTemplate(postgresConfigTemplate, postgresConfigModel, postgresConfigPath); err != nil {
		return render.WrapError("postgres config", err)
	}

	return nil
}

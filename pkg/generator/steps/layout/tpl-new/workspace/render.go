package workspace

import (
	_ "embed"
	"path/filepath"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed dotmify/gitignore.tpl
var gitignoreTemplate string

func Render(ctx *gencontext.GenContext) error {
	cacheGitignoreAbsPath := filepath.Join(ctx.GetWorkspace().GetCacheDirectory(), ".gitignore")
	if err := render.RenderTemplate(gitignoreTemplate, struct{}{}, cacheGitignoreAbsPath); err != nil {
		return err
	}
	return nil
}

package workspace

import (
	"github.com/chebykinn/mify/internal/mify/core"
)

func transformPath(context interface{}, path string) (string, error) {
	return path, nil
}

func RenderTemplateTree(context Context) error {
	params := core.RenderParams{
		TemplatesPath:   "tpl/workspace",
		TargetPath:      context.BasePath,
		PathTransformer: transformPath,
	}
	return core.RenderTemplateTree(context, params)
}

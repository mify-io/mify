package processors

import (
	"fmt"
	"path/filepath"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
)

type jsPostProcessor struct {
}

func newJsProcessor() *jsPostProcessor {
	return &jsPostProcessor{}
}

func (p *jsPostProcessor) GetServerGeneratorConfig(ctx *gencontext.GenContext) (GeneratorConfig, error) {
	basePath := ctx.GetWorkspace().BasePath
	targetPath, err := ctx.GetWorkspace().GetServiceDirectoryRelPath(
		ctx.GetServiceName(), ctx.MustGetMifySchema().Language, ctx.MustGetMifySchema().Template)
	if err != nil {
		return GeneratorConfig{}, err
	}
	generatedPath := filepath.Join(basePath, targetPath, "generated")
	return GeneratorConfig{
		TargetPath:  generatedPath,
		PackageName: SERVER_PACKAGE_NAME,
	}, nil
}

func (p *jsPostProcessor) GetClientGeneratorConfig(ctx *gencontext.GenContext, clientName string) (GeneratorConfig, error) {
	basePath := ctx.GetWorkspace().BasePath
	targetPath, err := ctx.GetWorkspace().GetServiceDirectoryRelPath(
		ctx.GetServiceName(), ctx.MustGetMifySchema().Language, ctx.MustGetMifySchema().Template)
	if err != nil {
		return GeneratorConfig{}, err
	}
	generatedPath := filepath.Join(basePath, targetPath, "generated", "api", "clients", clientName)
	packageName := endpoints.SanitizeServiceName(clientName) + "_client"
	return GeneratorConfig{
		TargetPath:  generatedPath,
		PackageName: packageName,
	}, nil
}

func (p *jsPostProcessor) ProcessServer(ctx *gencontext.GenContext) error {
	return nil
}

func (p *jsPostProcessor) ProcessClient(ctx *gencontext.GenContext, clientName string) error {
	return nil
}

func (p *jsPostProcessor) PopulateServerHandlers(ctx *gencontext.GenContext, paths []string) error {
	return fmt.Errorf("server handlers are not yet supported for js generator")
}

func (p *jsPostProcessor) Format(ctx *gencontext.GenContext) error {
	return nil
}

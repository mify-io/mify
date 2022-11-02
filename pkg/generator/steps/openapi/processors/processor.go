package processors

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

const (
	SERVER_PACKAGE_NAME = "openapi"
)

type GeneratorConfig struct {
	TargetPath string
	PackageName string
}

type PostProcessor interface {
	GetServerGeneratorConfig(ctx *gencontext.GenContext) (GeneratorConfig, error)
	GetClientGeneratorConfig(ctx *gencontext.GenContext, clientName string) (GeneratorConfig, error)
	ProcessServer(ctx *gencontext.GenContext) error
	ProcessClient(ctx *gencontext.GenContext, clientName string) error
	PopulateServerHandlers(ctx *gencontext.GenContext, paths []string) error
	Format(ctx *gencontext.GenContext) error
}

func NewPostProcessor(lang mifyconfig.ServiceLanguage) (PostProcessor, error) {
	switch lang {
	case mifyconfig.ServiceLanguagePython:
		return newPythonProcessor(), nil
	case mifyconfig.ServiceLanguageGo:
		return newGoProcessor(), nil
	case mifyconfig.ServiceLanguageJs:
		return newJsProcessor(), nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
}

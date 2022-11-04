package configuration

import (
	_ "embed"

	"github.com/mify-io/mify/internal/mify/generator/components/configuration/golang"
	"github.com/mify-io/mify/internal/mify/generator/components/configuration/python"
	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/workspace"
)

type ConfigurationComponent struct {
}

func NewConfigurationComponent() ConfigurationComponent {
	return ConfigurationComponent{}
}

func (s ConfigurationComponent) Name() string {
	return "configuration"
}

func (s ConfigurationComponent) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	switch ctx.GetServiceLanguage() {
	case workspace.Golang:
		return golang.Execute(ctx)
	case workspace.Python:
		return python.Execute(ctx)
	}

	return core.Done, nil
}

var _ core.Step = &ConfigurationComponent{}

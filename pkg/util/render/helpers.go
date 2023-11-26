package render

import (
	"fmt"
	"strings"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

func makeClientEnvName(serviceName string) string {
	sanitizedName := endpoints.SanitizeServiceName(serviceName)
	return strings.ToUpper(sanitizedName) + "_CLIENT_ENDPOINT"
}

type clientModel struct {
	ServiceName     string
	EndpointEnvName string
}

type DefaultModel struct {
	ServiceName string
	TplHeader string
	Clients []clientModel
	Workspace *WorkspaceModel
	Service *ServiceModel

	serviceLanguage mifyconfig.ServiceLanguage
}

type ModelWrapper[T any] struct {
	DefaultModel
	Model T
}

func newClientModel(clientName string) clientModel {
	return clientModel{
		ServiceName:     endpoints.SanitizeServiceName(clientName),
		EndpointEnvName: makeClientEnvName(clientName),
	}
}

func NewDefaultModel(ctx *gencontext.GenContext) DefaultModel {
	clients := []clientModel{}
	for clientName := range ctx.MustGetMifySchema().OpenAPI.Clients {
		clients = append(clients, newClientModel(clientName))
	}
	return DefaultModel{
		ServiceName: ctx.GetServiceName(),
		Clients:     clients,
		TplHeader: getTplHeader(ctx.MustGetMifySchema().Language),
		Workspace: NewWorkspaceModel(ctx),
		Service: NewServiceModel(ctx),
		serviceLanguage: ctx.MustGetMifySchema().Language,
	}
}

func NewModel[T any](ctx *gencontext.GenContext, model T) ModelWrapper[T] {
	return ModelWrapper[T]{
		NewDefaultModel(ctx),
		model,
	}
}

func getTplHeader(lang mifyconfig.ServiceLanguage) string {
	switch(lang) {
	case mifyconfig.ServiceLanguageGo, mifyconfig.ServiceLanguageJs:
		return "// THIS FILE IS AUTOGENERATED, DO NOT EDIT\n// Generated by mify"
	case mifyconfig.ServiceLanguagePython:
		return "# THIS FILE IS AUTOGENERATED, DO NOT EDIT\n# Generated by mify"
	}
	panic(fmt.Sprintf("unknown language: %s", lang))
}

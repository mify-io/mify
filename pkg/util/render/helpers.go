package render

import (
	"strings"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
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
	// TplHeader string
	Clients []clientModel
}

func newClientModel(clientName string) clientModel {
	return clientModel{
		ServiceName:     endpoints.SanitizeServiceName(clientName),
		EndpointEnvName: makeClientEnvName(clientName),
	}
}

func NewDefaultModel(ctx *gencontext.GenContext) DefaultModel {
	clients := []clientModel{}
	for clientName := range ctx.GetMifySchema().OpenAPI.Clients {
		clients = append(clients, newClientModel(clientName))
	}
	return DefaultModel{
		ServiceName: ctx.GetServiceName(),
		Clients:     clients,
	}
}

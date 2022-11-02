package core

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
)

type clientModel struct {
	ApiClientName     string
	ConfigurationName string
	PropertyName      string
	ImportPath        string
}

type clientsModel struct {
	TplHeader   string
	ServiceName string

	// MetricsIncludePath string
	Clients []clientModel
}

func newClientModel(serviceName string, clientName string) clientModel {
	clientNameSanitized := endpoints.SanitizeServiceName(clientName)
	clientNameCamelCase := endpoints.SnakeCaseToCamelCase(clientNameSanitized, true)

	return clientModel{
		ApiClientName:     clientNameCamelCase + "ApiClient",
		ConfigurationName: clientNameCamelCase + "Configuration",
		PropertyName:      clientNameSanitized + "_client",
		ImportPath:        fmt.Sprintf("%s.generated.openapi.clients.%s", serviceName, clientNameSanitized),
	}
}

func newClientsModel(ctx *gencontext.GenContext) clientsModel {
	serviceName := ctx.GetMifySchema().ServiceName
	clients := []clientModel{}
	for clientName := range ctx.GetMifySchema().OpenAPI.Clients {
		clients = append(clients, newClientModel(serviceName, clientName))
	}
	return clientsModel{
		TplHeader:   ctx.GetWorkspace().TplHeaderPy,
		ServiceName: ctx.GetMifySchema().ServiceName,
		Clients:     clients,
	}
}

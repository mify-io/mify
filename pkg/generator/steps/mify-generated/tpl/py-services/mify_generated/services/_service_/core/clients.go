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
	Clients []clientModel
}

func newClientModel(svcImport string, serviceName string, clientName string) clientModel {
	clientNameSanitized := endpoints.SanitizeServiceName(clientName)
	clientNameCamelCase := endpoints.SnakeCaseToCamelCase(clientNameSanitized, true)

	return clientModel{
		ApiClientName:     clientNameCamelCase + "ApiClient",
		ConfigurationName: clientNameCamelCase + "Configuration",
		PropertyName:      clientNameSanitized + "_client",
		ImportPath:        fmt.Sprintf("%s.openapi.clients.%s", svcImport, clientNameSanitized),
	}
}

func newClientsModel(ctx *gencontext.GenContext) clientsModel {
	serviceName := ctx.GetMifySchema().ServiceName
	clients := []clientModel{}
	mifyGen := ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema())
	svcImport := mifyGen.GetServicePackage()
	for clientName := range ctx.GetMifySchema().OpenAPI.Clients {
		clients = append(clients, newClientModel(svcImport, serviceName, clientName))
	}
	return clientsModel{
		Clients:     clients,
	}
}

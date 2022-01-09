package apigateway

import (
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/service/client"
	"github.com/chebykinn/mify/internal/mify/workspace"
)

func RegenerateHandlers(ctx *core.Context, workspace workspace.Context, publicApis PublicApis) error {
	for _, serviceName := range publicApis {
		if err := client.AddClient(ctx, workspace, ApiGatewayName, serviceName); err != nil {
			return err
		}
	}

	return nil
}

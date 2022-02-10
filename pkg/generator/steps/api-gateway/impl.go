package apigateway

import (
	"fmt"

	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/workspace/mutators"
	"github.com/mify-io/mify/pkg/workspace/mutators/client"
)

func execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	foundPublicApis := scanPublicApis(ctx)
	openapiUpdated, err := updateApiGatewayOpenapiSchema(ctx, foundPublicApis)
	if err != nil {
		return core.Done, fmt.Errorf("can't update openapi schema: %w", err)
	}

	clientsUpadted, err := updateApiGatewayClients(ctx, foundPublicApis)
	if err != nil {
		return core.Done, fmt.Errorf("can't update clients: %w", err)
	}

	if openapiUpdated || clientsUpadted {
		return core.RepeatAll, nil
	}

	return core.Done, nil
}

func updateApiGatewayClients(ctx *gencontext.GenContext, publicApis PublicApis) (bool, error) {
	mutCtx := mutators.NewMutatorContext(ctx.GetGoContext(), nil, ctx.GetWorkspace())
	updated := false

	for targetService := range publicApis {
		if ctx.GetMifySchema().OpenAPI.HasClient(targetService) {
			continue
		}

		updated = true
		if err := client.AddClient(mutCtx, ctx.GetServiceName(), targetService); err != nil {
			return updated, err
		}
	}

	for targetService := range ctx.GetMifySchema().OpenAPI.Clients {
		_, ok := publicApis[targetService]
		if ok {
			continue
		}

		updated = true
		if err := client.RemoveClient(mutCtx, ctx.GetServiceName(), targetService); err != nil {
			return updated, err
		}
	}

	return updated, nil
}

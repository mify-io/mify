package openapi

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

func isFullRegenerationNeeded(ctx *gencontext.GenContext) (bool, error) {
	meta, err := getGenerationMeta(ctx)
	if err != nil {
		return false, err
	}
	if meta.MifyVersion != ctx.GetMifyVersion() {
		return true, nil
	}
	if ctx.GetForceRegeneration() {
		return true, nil
	}
	return false, nil
}

func generateServiceOpenAPI(ctx *gencontext.GenContext) error {
	openapigen, err := NewOpenAPIGenerator(ctx)
	if err != nil {
		return err
	}

	forceRegeneration, err := isFullRegenerationNeeded(ctx)
	if err != nil {
		return err
	}
	if forceRegeneration {
		ctx.Logger.Infof("Will do full regeneration")
	}

	serverNeedsRegeneration, err := checkServerNeedsRegeneration(ctx, &openapigen, forceRegeneration)
	if err != nil {
		return err
	}

	if !serverNeedsRegeneration {
		ctx.Logger.Infof("Server side of '%s' service is actual. Skipping...", ctx.GetServiceName())
	}

	clientsDiff, err := calcClientsDiff(ctx, &openapigen, forceRegeneration)
	if err != nil {
		return err
	}

	if clientsDiff.Empty() {
		ctx.Logger.Infof("Clients included in service '%s' are actual. Skipping...", ctx.GetServiceName())
	}

	return doGeneration(ctx, &openapigen, serverNeedsRegeneration, clientsDiff)
}

func doGeneration(
	ctx *gencontext.GenContext,
	openapigen *OpenAPIGenerator,
	generateServer bool,
	clientsDiff clientsDiff) error {

	executePool := ctx.GetExecutePoolFactory().NewPool()

	executePool.EnqueExecution(func() error {
		if !generateServer {
			return nil
		}

		if err := tryPrepareOpenApi(ctx, openapigen); err != nil {
			return err
		}

		if err := generateServerSide(ctx, openapigen); err != nil {
			return fmt.Errorf("failed while generating server side: %w", err)
		}

		return nil
	})

	executePool.EnqueExecution(func() error {
		if clientsDiff.Empty() {
			return nil
		}

		if err := tryPrepareOpenApi(ctx, openapigen); err != nil {
			return err
		}

		if err := generateClients(ctx, openapigen, clientsDiff); err != nil {
			return fmt.Errorf("failed while generating clients: %w", err)
		}

		return nil
	})

	executePool.EnqueExecution(func() error {
		if !needGenerateClientsContext(ctx, clientsDiff) {
			return nil
		}

		if err := generateClientsContext(ctx); err != nil {
			return fmt.Errorf("failed while generating clients context: %w", err)
		}

		return nil
	})

	errs := executePool.WaitAll()
	if errs != nil {
		return errs[0]
	}

	if err := updateClientsList(ctx); err != nil {
		return fmt.Errorf("failed while updating mify schema: %w", err)
	}

	if err := writeGenerationMeta(ctx, GenerationMeta{MifyVersion: ctx.GetMifyVersion()}); err != nil {
		return fmt.Errorf("failed while updating mify schema: %w", err)
	}

	return nil
}

func tryPrepareOpenApi(ctx *gencontext.GenContext, openapigen *OpenAPIGenerator) error {
	err := openapigen.PrepareSync(ctx)
	if err != nil {
		return fmt.Errorf("failed while preparing openapi: %w", err)
	}

	return nil
}

func generateServerSide(ctx *gencontext.GenContext, openAPIGenerator *OpenAPIGenerator) error {
	ctx.Logger.Infof("Generating server side of service '%s'", ctx.GetServiceName())

	if err := openAPIGenerator.GenerateServer(ctx); err != nil {
		return err
	}
	return nil
}

func checkServerNeedsRegeneration(ctx *gencontext.GenContext, openapigen *OpenAPIGenerator, force bool) (bool, error) {
	wrapError := func(err error) error {
		return fmt.Errorf("can't check if server needs regeneration: %w", err)
	}

	if !hasServerApiSchema(ctx) {
		return false, nil
	}

	schemaDirPath := ctx.GetWorkspace().GetApiSchemaDirRelPath(ctx.GetServiceName())
	needGenerateServer, err := openapigen.NeedGenerateServer(ctx, schemaDirPath)
	if err != nil {
		return false, wrapError(err)
	}

	return needGenerateServer || force, nil
}

func generateClients(ctx *gencontext.GenContext, openapigen *OpenAPIGenerator, clientsDiff clientsDiff) error {
	ctx.Logger.Infof("Generating clients inside service '%s'", ctx.GetServiceName())

	executePool := ctx.GetExecutePoolFactory().NewPool()

	for clientName := range clientsDiff.removed {
		executePool.EnqueExecution(func() error {
			ctx.Logger.Debugf("Removing client '%s' from service '%s' ...", clientName, ctx.GetServiceName())

			err := openapigen.RemoveClient(ctx, clientName)
			if err != nil {
				return err
			}

			return nil
		})
	}

	for clientName := range clientsDiff.added {
		executePool.EnqueExecution(func() error {
			ctx.Logger.Debugf("Adding client '%s' to service '%s' ...", clientName, ctx.GetServiceName())

			if err := openapigen.GenerateClient(ctx, clientName); err != nil {
				return fmt.Errorf("failed to generate client for: %s: %w", clientName, err)
			}

			return nil
		})
	}

	for clientName := range clientsDiff.schemaChanged {
		executePool.EnqueExecution(func() error {
			ctx.Logger.Debugf("Regenerating client '%s' in service '%s' ...", clientName, ctx.GetServiceName())

			if err := openapigen.GenerateClient(ctx, clientName); err != nil {
				return fmt.Errorf("failed to generate client for: %s: %w", clientName, err)
			}

			return nil
		})
	}

	errs := executePool.WaitAll()
	if errs != nil {
		return errs[0]
	}

	return nil
}

// Some services could don't have scheme (f.e. frontend)
func hasServerApiSchema(ctx *gencontext.GenContext) bool {
	schemas := ctx.GetSchemaCtx().MustGetServiceSchemas(ctx.GetServiceName())
	return len(schemas.GetOpenapi()) > 0
}

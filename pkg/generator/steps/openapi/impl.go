package openapi

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/containerd/containerd/log"
)

func generateServiceOpenAPI(ctx *gencontext.GenContext) error {
	openapigen := NewOpenAPIGenerator(ctx)

	serverNeedsRegeneration, err := checkServerNeedsRegeneration(ctx, openapigen)
	if err != nil {
		return err
	}

	if !serverNeedsRegeneration {
		ctx.Logger.Infof("Server side of '%s' service is actual. Skipping...", ctx.GetServiceName())
	}

	clientsDiff, err := calcClientsDiff(ctx, &openapigen)
	if err != nil {
		return err
	}

	if clientsDiff.Empty() {
		ctx.Logger.Infof("Clients included in service '%s' are actual. Skipping...", ctx.GetServiceName())
	}

	return doGeneration(ctx, openapigen, serverNeedsRegeneration, clientsDiff)
}

func doGeneration(
	ctx *gencontext.GenContext,
	openapigen OpenAPIGenerator,
	generateServer bool,
	clientsDiff clientsDiff) error {

	wrapError := func(err error) error {
		return fmt.Errorf("error during generation: %w", err)
	}

	// TODO: server + clients parallelization

	if generateServer || !clientsDiff.Empty() {
		err := openapigen.Prepare(ctx)
		if err != nil {
			return wrapError(err)
		}

		if generateServer {
			err := generateServerSide(ctx, &openapigen)
			if err != nil {
				return wrapError(err)
			}
		}

		if !clientsDiff.Empty() {
			err = generateClients(ctx, &openapigen, clientsDiff)
			if err != nil {
				return wrapError(err)
			}
		}
	}

	if needGenerateClientsContext(ctx, clientsDiff) {
		err := generateClientsContext(ctx)
		if err != nil {
			return wrapError(err)
		}

		err = updateClientsList(ctx)
		if err != nil {
			return wrapError(err)
		}
	}

	return nil
}

func generateServerSide(ctx *gencontext.GenContext, openAPIGenerator *OpenAPIGenerator) error {
	ctx.Logger.Infof("Generating server side of service '%s'", ctx.GetServiceName())

	targetDir, err := getAPIServicePathByLang(ctx.MustGetMifySchema().Language, ctx.GetServiceName())
	if err != nil {
		return err
	}

	if err := openAPIGenerator.GenerateServer(ctx, targetDir); err != nil {
		return err
	}
	return nil
}

func checkServerNeedsRegeneration(ctx *gencontext.GenContext, openapigen OpenAPIGenerator) (bool, error) {
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

	return needGenerateServer, nil
}

func generateClients(ctx *gencontext.GenContext, openapigen *OpenAPIGenerator, clientsDiff clientsDiff) error {
	ctx.Logger.Infof("Generating clients inside service '%s'", ctx.GetServiceName())

	targetDir, err := getAPIServicePathByLang(ctx.MustGetMifySchema().Language, ctx.GetServiceName())
	if err != nil {
		return err
	}

	// TODO: parallel

	for clientName := range clientsDiff.removed {
		log.L.Trace("Removing client '%s' from service '%s' ...", clientName, ctx.GetServiceName())

		err := openapigen.RemoveClient(ctx, clientName, targetDir)
		if err != nil {
			return err
		}
	}

	for clientName := range clientsDiff.added {
		log.L.Trace("Adding client '%s' to service '%s' ...", clientName, ctx.GetServiceName())

		if err := openapigen.GenerateClient(ctx, clientName, targetDir); err != nil {
			return fmt.Errorf("failed to generate client for: %s: %w", clientName, err)
		}
	}

	for clientName := range clientsDiff.schemaChanged {
		log.L.Trace("Regenerating client '%s' in service '%s' ...", clientName, ctx.GetServiceName())

		if err := openapigen.GenerateClient(ctx, clientName, targetDir); err != nil {
			return fmt.Errorf("failed to generate client for: %s: %w", clientName, err)
		}
	}

	return nil
}

func getAPIServicePathByLang(language mifyconfig.ServiceLanguage, serviceName string) (string, error) {
	switch language {
	case mifyconfig.ServiceLanguageGo:
		return mifyconfig.GoServicesRoot + "/internal/" + serviceName, nil
	case mifyconfig.ServiceLanguageJs:
		return mifyconfig.JsServicesRoot + "/" + serviceName, nil
	}
	return "", fmt.Errorf("unknown language: %s", language)
}

// Some services could don't have scheme (f.e. frontend)
func hasServerApiSchema(ctx *gencontext.GenContext) bool {
	schemas := ctx.GetSchemaCtx().MustGetServiceSchemas(ctx.GetServiceName())
	return len(schemas.GetOpenapi()) > 0
}

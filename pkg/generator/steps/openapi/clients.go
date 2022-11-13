package openapi

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
	"github.com/mify-io/mify/pkg/generator/steps/openapi/tpl"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed tpl/go_clients.go.tpl
var goClientsTemplate string

//go:embed tpl/js_clients.js.tpl
var jsClientsTemplate string

// Update only when new services added or removed or context is not generated yet
func needGenerateClientsContext(ctx *gencontext.GenContext, clientsDiff clientsDiff) bool {
	pathToClientsContext, err := getAbsPathToClientsContext(ctx)
	if err != nil {
		return false
	}
	_, err = os.Stat(pathToClientsContext)
	return os.IsNotExist(err) || len(clientsDiff.added) > 0 || len(clientsDiff.removed) > 0
}

func getAbsPathToClientsContext(ctx *gencontext.GenContext) (string, error) {
	switch ctx.MustGetMifySchema().Language {
	case mifyconfig.ServiceLanguageGo:
		generatedDirPath := ctx.GetWorkspace().GetGeneratedAbsPath(ctx.GetServiceName())
		return path.Join(generatedDirPath, "core", "clients.go"), nil
	case mifyconfig.ServiceLanguageJs:
		generatedDirPath := ctx.GetWorkspace().GetJsGeneratedAbsPath(ctx.GetServiceName())
		return path.Join(generatedDirPath, "core", "clients.js"), nil
	}

	return "", errors.New("unknown or not supported language")
}

// Generate struct which will be included in service context (generated part of service)
func generateClientsContext(ctx *gencontext.GenContext) error {
	ctx.Logger.Infof("Generating clients context in service '%s'", ctx.GetServiceName())

	path, err := getAbsPathToClientsContext(ctx)
	if err != nil {
		return err
	}

	switch ctx.MustGetMifySchema().Language {
	case mifyconfig.ServiceLanguageGo:
		clientsModel, err := makeGoClientsModel(ctx)
		if err != nil {
			return err
		}

		if err := render.RenderTemplate(goClientsTemplate, clientsModel, path); err != nil {
			return err
		}
	case mifyconfig.ServiceLanguageJs:
		clientsModel, err := makeJsClientsModel(ctx)
		if err != nil {
			return err
		}

		if err := render.RenderTemplate(jsClientsTemplate, clientsModel, path); err != nil {
			return err
		}
	}

	return nil
}

func makeGoClientsModel(ctx *gencontext.GenContext) (tpl.GoClientsModel, error) {
	targetServices := ctx.MustGetMifySchema().OpenAPI.Clients
	clientsList := make([]tpl.GoClientModel, 0, len(targetServices))
	for targetServiceName := range targetServices {
		targetServiceSchemas := ctx.GetSchemaCtx().MustGetServiceSchemas(targetServiceName)
		if len(targetServiceSchemas.GetOpenapi()) == 0 {
			return tpl.GoClientsModel{}, fmt.Errorf("schema of '%s' wasn't found while generating client in '%s'", targetServiceName, ctx.GetServiceName())
		}

		svcNameSanitized := endpoints.SanitizeServiceName(targetServiceName)
		packageName := svcNameSanitized + "_client"
		fieldName := endpoints.SnakeCaseToCamelCase(svcNameSanitized, false)
		methodName := endpoints.SnakeCaseToCamelCase(svcNameSanitized, true)
		clientsList = append(clientsList, tpl.NewGoClientModel(
			targetServiceName,
			packageName,
			fieldName,
			methodName,
			fmt.Sprintf("%s/internal/%s/generated/api/clients/%s",
				ctx.GetWorkspace().GetGoModule(),
				ctx.GetServiceName(),
				targetServiceName),
		))
	}
	sort.Slice(clientsList, func(i, j int) bool {
		return clientsList[i].ClientName < clientsList[j].ClientName
	})

	metricsIncludePath := ""
	if len(clientsList) != 0 {
		metricsIncludePath = fmt.Sprintf("%s/internal/pkg/generated/metrics", ctx.GetWorkspace().GetGoModule())
	}

	return tpl.NewGoClientsModel(
		ctx.GetWorkspace().TplHeader,
		ctx.GetWorkspace().GoRoot,
		ctx.GetServiceName(),
		metricsIncludePath,
		clientsList), nil
}

func makeJsClientsModel(ctx *gencontext.GenContext) (tpl.JsClientsModel, error) {
	targetServices := ctx.MustGetMifySchema().OpenAPI.Clients
	clientsList := make([]tpl.JsClientModel, 0, len(targetServices))
	for targetServiceName := range targetServices {
		targetServiceSchemas := ctx.GetSchemaCtx().MustGetServiceSchemas(targetServiceName)
		if len(targetServiceSchemas.GetOpenapi()) == 0 {
			return tpl.JsClientsModel{}, fmt.Errorf("schema of '%s' wasn't found while generating client in '%s'", targetServiceName, ctx.GetServiceName())
		}

		methodName := endpoints.SnakeCaseToCamelCase(endpoints.SanitizeServiceName(targetServiceName), false)
		className := endpoints.SnakeCaseToCamelCase(endpoints.SanitizeServiceName(targetServiceName), true)
		clientsList = append(clientsList, tpl.NewJsClientModel(
			targetServiceName,
			className,
			methodName,
		))
	}
	sort.Slice(clientsList, func(i, j int) bool {
		return clientsList[i].ClientName < clientsList[j].ClientName
	})

	return tpl.NewJsClientsModel(
		ctx.GetWorkspace().TplHeader,
		ctx.GetServiceName(),
		clientsList), nil
}

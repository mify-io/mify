package processors

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
)

type jsPostProcessor struct {
}

func newJsProcessor() *jsPostProcessor {
	return &jsPostProcessor{}
}

func (p *jsPostProcessor) GetServerGeneratorConfig(ctx *gencontext.GenContext) (GeneratorConfig, error) {
	basePath := ctx.GetWorkspace().BasePath
	targetPath, err := ctx.GetWorkspace().GetServiceDirectoryRelPath(
		ctx.GetServiceName(), ctx.MustGetMifySchema().Language, ctx.MustGetMifySchema().Template)
	if err != nil {
		return GeneratorConfig{}, err
	}
	generatedPath := filepath.Join(basePath, targetPath, "generated")
	return GeneratorConfig{
		TargetPath:  generatedPath,
		PackageName: SERVER_PACKAGE_NAME,
	}, nil
}

func (p *jsPostProcessor) GetClientGeneratorConfig(ctx *gencontext.GenContext, clientName string) (GeneratorConfig, error) {
	basePath := ctx.GetWorkspace().BasePath
	targetPath, err := ctx.GetWorkspace().GetServiceDirectoryRelPath(
		ctx.GetServiceName(), ctx.MustGetMifySchema().Language, ctx.MustGetMifySchema().Template)
	if err != nil {
		return GeneratorConfig{}, err
	}
	generatedPath := filepath.Join(basePath, targetPath, "generated", "api", "clients", clientName)
	packageName := endpoints.SanitizeServiceName(clientName) + "_client"
	return GeneratorConfig{
		TargetPath:  generatedPath,
		PackageName: packageName,
	}, nil
}

func (p *jsPostProcessor) ProcessServer(ctx *gencontext.GenContext) error {
	return nil
}

func (p *jsPostProcessor) ProcessClient(ctx *gencontext.GenContext, clientName string) error {
	return nil
}

func (p *jsPostProcessor) PopulateServerHandlers(ctx *gencontext.GenContext, paths []string) error {
	handlersTargetDir, err := ctx.GetWorkspace().GetServiceGeneratedAPIRelPath(
		ctx.GetServiceName(), ctx.MustGetMifySchema().Language)
	if err != nil {
		return err
	}

	pathsSet := map[string]string{}
	for _, path := range paths {
		ctx.Logger.Infof("pre path: %s", path)
		path = strings.ReplaceAll(path, "{", "")
		path = strings.ReplaceAll(path, "}", "")
		pathsSet[p.toAPIFilename(path)] = path
	}

	ctx.Logger.Infof("paths: %v", pathsSet)

	err = p.moveOutHandlers(ctx, pathsSet, handlersTargetDir)
	if err != nil {
		return err
	}

	err = p.fixImportsForControllers(ctx, handlersTargetDir, pathsSet)
	if err != nil {
		return err
	}

	return err
}

func (p *jsPostProcessor) moveServicesIndexJs(servicesPath string, handlersTargetPath string) error {
	input, err := os.ReadFile(filepath.Join(servicesPath, "index.js"))
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(handlersTargetPath, "index.js"), input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *jsPostProcessor) moveOutHandlers(ctx *gencontext.GenContext, pathsSet map[string]string, handlersTargetDir string) error {
	handlerGlob := "*Service.js"
	handlerSuffix := "Service.js"
	targetServiceName := "service.js"

	servicesPath := filepath.Join(ctx.GetWorkspace().BasePath, handlersTargetDir, "generated", "services")

	services, err := filepath.Glob(filepath.Join(servicesPath, handlerGlob))
	if err != nil {
		return err
	}
	ctx.Logger.Infof("services: %v", services)
	if len(services) == 0 {
		ctx.Logger.Infof("no handlers to move")
		return nil
	}

	handlersTargetPath := filepath.Join(ctx.GetWorkspace().BasePath, handlersTargetDir, "handlers")
	// err = sanitizeJsServerHandlersImports(ctx, apiPath)
	if err != nil {
		return err
	}

	for _, service := range services {
		serviceFileName := filepath.Base(service)
		serviceFileName = strings.TrimSuffix(serviceFileName, handlerSuffix)

		path, ok := pathsSet[serviceFileName]
		if !ok {
			if serviceFileName == "" {
				path = "" // Move Service.js (misc file) to root of handlers
			} else {
				return fmt.Errorf("failed to find path for service file: %s", serviceFileName)
			}
		}

		ctx.Logger.Infof("processing handler for: %v", path)
		targetFile := filepath.Join(handlersTargetPath, path, targetServiceName)
		if _, err := os.Stat(targetFile); err == nil {
			ctx.Logger.Infof("skipping existing handler for: %v", path)
			continue
		}
		if err := os.MkdirAll(filepath.Join(handlersTargetPath, path), 0755); err != nil {
			return err
		}
		if err := p.moveServiceToHandler(ctx, service, handlersTargetPath, targetFile); err != nil {
			return err
		}
		ctx.Logger.Infof("created handler for: %v", path)
	}

	err = p.moveServicesIndexJs(servicesPath, handlersTargetPath)
	if err != nil {
		return err
	}

	err = os.RemoveAll(servicesPath)
	if err != nil {
		return fmt.Errorf("can't remove old services directory: %w", err)
	}

	return nil
}

func (p *jsPostProcessor) fixImportsForControllers(ctx *gencontext.GenContext, handlersTargetDir string, pathsSet map[string]string) error {
	controllerGlob := "*Controller.js"
	controllerSuffix := "Controller.js"
	controllersPath := filepath.Join(ctx.GetWorkspace().BasePath, handlersTargetDir, "generated", "controllers")

	controllers, err := filepath.Glob(filepath.Join(controllersPath, controllerGlob))
	if err != nil {
		return err
	}
	ctx.Logger.Infof("controllers: %v", controllers)
	if len(controllers) == 0 {
		ctx.Logger.Infof("no controllers to fix import")
		return nil
	}

	for _, controller := range controllers {
		controllerFileName := filepath.Base(controller)
		controllerFileName = strings.TrimSuffix(controllerFileName, controllerSuffix)
		if controllerFileName == "" {
			continue // Ignore misc files
		}

		path, ok := pathsSet[controllerFileName]
		if !ok {
			return fmt.Errorf("failed to find path for controller file: %s", controllerFileName)
		}

		input, err := os.ReadFile(filepath.Join(controller))
		if err != nil {
			return err
		}

		m1 := regexp.MustCompile(`const service = require\(['"].*['"]\);`)
		data := m1.ReplaceAllString(
			string(input),
			fmt.Sprintf("const service = require('../../handlers/%s/service.js');", strings.TrimPrefix(path, "/")))

		err = os.WriteFile(controller, []byte(data), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *jsPostProcessor) toAPIFilename(name string) string {
	// NOTE: openapi-generator transforms tag to camelCase, we don't do that here
	// we just remove slashes from path and then use openapi-generator logic
	// to convert this path to filename.
	api := strings.TrimPrefix(name, "/")
	api = strings.ReplaceAll(api, "/", "_")
	// replace - with _ e.g. created-at => created_at
	api = strings.ReplaceAll(api, "-", "_")
	return endpoints.SnakeCaseToCamelCase(api, true)
}

func (p *jsPostProcessor) Format(ctx *gencontext.GenContext) error {
	return nil
}

func (p *jsPostProcessor) moveServiceToHandler(ctx *gencontext.GenContext, serviceFile string, handlersTargetPath string, targetFile string) error {
	input, err := os.ReadFile(serviceFile)
	if err != nil {
		return err
	}

	targetRelPath := strings.TrimPrefix(targetFile, handlersTargetPath)
	targetRelPath = strings.TrimPrefix(targetRelPath, "/")

	depth := strings.Count(targetRelPath, "/")

	data := strings.ReplaceAll(string(input), "require('./Service')", "require('./service')")

	m1 := regexp.MustCompile(`require\(['"](.*)['"]\);`)
	data = m1.ReplaceAllStringFunc(data, func(s string) string {
		return strings.Replace(s, "./", p.constructImportPrefix(depth), 1)
	})

	err = os.WriteFile(targetFile, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *jsPostProcessor) constructImportPrefix(depth int) string {
	if depth == 1 {
		return "./"
	}

	var sb strings.Builder
	for i := 0; i < depth; i++ {
		sb.WriteString("../")
	}

	return sb.String()
}

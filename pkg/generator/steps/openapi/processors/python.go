package processors

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
)

type pythonPostProcessor struct {
}

func newPythonProcessor() *pythonPostProcessor {
	return &pythonPostProcessor{}
}

func (p *pythonPostProcessor) GetServerGeneratorConfig(ctx *gencontext.GenContext) (GeneratorConfig, error) {
	generatedPath := ctx.GetWorkspace().GetPythonServicesAbsPath()
	return GeneratorConfig{
		TargetPath:  generatedPath,
		PackageName: fmt.Sprintf("%s.%s.%s", endpoints.SanitizeServiceName(ctx.GetServiceName()), "generated", SERVER_PACKAGE_NAME),
	}, nil
}

func (p *pythonPostProcessor) GetClientGeneratorConfig(ctx *gencontext.GenContext, clientName string) (GeneratorConfig, error) {
	generatedPath := ctx.GetWorkspace().GetPythonServicesAbsPath()
	return GeneratorConfig{
		TargetPath:  generatedPath,
		PackageName: fmt.Sprintf("%s.%s.%s.%s.%s", endpoints.SanitizeServiceName(ctx.GetServiceName()), "generated", SERVER_PACKAGE_NAME, "clients", clientName),
	}, nil
}

func (p *pythonPostProcessor) ProcessServer(ctx *gencontext.GenContext) error {
	generatedPath := ctx.GetWorkspace().GetPythonServicesAbsPath()
	ignoreFilePath := filepath.Join(generatedPath, ".openapi-generator-ignore")
	openAPIGeneratorDir := filepath.Join(generatedPath, ".openapi-generator")
	if err := os.Remove(ignoreFilePath); err != nil {
		return err
	}

	if err := os.RemoveAll(openAPIGeneratorDir); err != nil {
		return err
	}

	if err := processControllers(ctx); err != nil {
		return err
	}

	return nil
}

func (p *pythonPostProcessor) ProcessClient(ctx *gencontext.GenContext, clientName string) error {
	generatedPath := filepath.Join(
		ctx.GetWorkspace().GetPythonServicesAbsPath(),
		endpoints.SanitizeServiceName(ctx.GetServiceName()),
		"generated", SERVER_PACKAGE_NAME, "clients", clientName)
	// TODO: remove after https://github.com/OpenAPITools/openapi-generator/issues/13648 is fixed
	err := filepath.WalkDir(generatedPath, func(p string, d fs.DirEntry, e error) error {
		if d == nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if e != nil {
			return e
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(bytes.NewReader(data))
		lines := make([]string, 0)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		buf := bytes.NewBufferString("")
		w := bufio.NewWriter(buf)
		for i, line := range lines {
			if strings.HasPrefix(line, "from") {
				line = strings.ReplaceAll(line, "/", ".")
			}
			if _, err := w.WriteString(line); err != nil {
				return err
			}
			if i+1 < len(lines) {
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}
		}
		err = w.Flush()
		if err != nil {
			return err
		}

		err = os.WriteFile(p, buf.Bytes(), 0666)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (p *pythonPostProcessor) PopulateServerHandlers(ctx *gencontext.GenContext, paths []string) error {
	var handlerGlob, handlerSuffix, targetServiceName string
	handlerGlob = "controllers/*_controller_service.py"
	handlerSuffix = "_controller_service.py"
	targetServiceName = "service.py"
	targetDir, err := ctx.GetWorkspace().GetServiceGeneratedAPIRelPath(
		ctx.GetServiceName(), ctx.MustGetMifySchema().Language)
	if err != nil {
		return err
	}
	generatedPath := filepath.Join(ctx.GetWorkspace().BasePath, targetDir, "generated")
	apiPath := filepath.Join(generatedPath, "openapi")
	handlersPath := filepath.Join(ctx.GetWorkspace().BasePath, targetDir, "handlers")
	services, err := filepath.Glob(filepath.Join(apiPath, handlerGlob))
	if err != nil {
		return err
	}
	ctx.Logger.Infof("services: %v", services)
	if len(services) == 0 {
		ctx.Logger.Infof("no handlers to move")
		return nil
	}
	pathsSet := map[string]string{}
	for _, path := range paths {
		path = strings.ReplaceAll(path, "{", "")
		path = strings.ReplaceAll(path, "}", "")
		path = strings.ReplaceAll(path, "-", "_")
		pathsSet[p.toAPIFilename(path)] = path
	}
	ctx.Logger.Infof("paths: %v", pathsSet)
	for _, service := range services {
		serviceFileName := filepath.Base(service)
		serviceFileName = strings.TrimSuffix(serviceFileName, handlerSuffix)
		path, ok := pathsSet[serviceFileName]
		if !ok {
			return fmt.Errorf("failed to find path for service file: %s", serviceFileName)
		}

		ctx.Logger.Infof("processing handler for: %v", path)
		targetFile := filepath.Join(handlersPath, path, targetServiceName)
		defer func(svc string) {
			if err := os.Remove(svc); err != nil {
				ctx.Logger.Infof("failed to remove service file: %s: %s", svc, err)
				return
			}
			ctx.Logger.Infof("cleaned generated service file: %s", svc)
		}(service)

		if _, err := os.Stat(targetFile); err == nil {
			ctx.Logger.Infof("skipping existing handler for: %v", path)
			continue
		}
		if err := os.MkdirAll(filepath.Join(handlersPath, path), 0755); err != nil {
			return err
		}
		if err := p.createServerHandlersFile(ctx, service, targetFile); err != nil {
			return err
		}
		ctx.Logger.Infof("created handler for: %v", path)
	}
	return nil
}
func (p *pythonPostProcessor) Format(ctx *gencontext.GenContext) error {
	return nil
}

func (p *pythonPostProcessor) toAPIFilename(name string) string {
	// NOTE: openapi-generator transforms tag to camelCase, we don't do that here
	// we just remove slashes from path and then use openapi-generator logic
	// to convert this path to filename.
	api := strings.TrimPrefix(name, "/")
	api = strings.ReplaceAll(api, "/", "_")
	// replace - with _ e.g. created-at => created_at
	api = strings.ReplaceAll(api, "-", "_")
	// // e.g. PetApi.go => pet_api.go
	api = underscore(api)
	return api
}

func (p *pythonPostProcessor) createServerHandlersFile(ctx *gencontext.GenContext, serviceFile string, targetFile string) error {
	data, err := os.ReadFile(serviceFile)
	if err != nil {
		return err
	}
	err = os.WriteFile(targetFile, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func processControllers(ctx *gencontext.GenContext) error {
	controllersPath := filepath.Join(ctx.GetWorkspace().BasePath,
		ctx.GetWorkspace().GetPythonServiceGeneratedOpenAPIRelPath(ctx.GetServiceName()), "controllers")
	controllers, err := filepath.Glob(filepath.Join(controllersPath, "*_controller.py"))
	if err != nil {
		return err
	}
	for _, controllerPath := range controllers {
		err = processController(controllerPath)
		if err != nil {
			return fmt.Errorf("failed to process controller: %s: %w",
				controllerPath, err)
		}
	}
	return nil
}

// TODO: remove copypasta - either make common line by line
// iterator, or just use go templates for everything
func processController(filename string) error {
	const (
		sectionStart = "# import_start"
		sectionEnd   = "# import_end"
	)

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	buf := bytes.NewBufferString("")
	w := bufio.NewWriter(buf)

	isSectionStart := false
	for i, line := range lines {
		if line == sectionStart {
			isSectionStart = true
			continue
		}
		if line == sectionEnd {
			isSectionStart = false
			continue
		}
		if !isSectionStart {
			if _, err := w.WriteString(line); err != nil {
				return err
			}

			if i+1 < len(lines) && lines[i+1] != sectionStart {
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}
			continue
		}
		line = strings.ReplaceAll(line, "/", ".")
		line = strings.ReplaceAll(line, "{", "")
		line = strings.ReplaceAll(line, "}", "")
		if _, err := w.WriteString(line); err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, buf.Bytes(), 0666)
	if err != nil {
		return err
	}

	return nil
}

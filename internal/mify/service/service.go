package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chebykinn/mify/internal/mify/config"
)

var (
	apiTemplate = `
	openapi: "3.0.0"
	info:
	  version: 1.0.0
	  title: %s
	  description: Service description
	  contact:
	    name: Maintainer name
	    email: Maintainer email
	    url: url
	servers:
	  - url: %s
	paths:
	  /path/to/api:
	    get:
	      summary: sample handler
	      operationId: theOperationId
	      responses:
	        '200':
	          description: OK
	          content:
	            application/json:
	              schema:
	                type: object
	`
)

func CreateService(wspConext workspace.Context, name string) error {
	fmt.Printf("creating service %s\n", name)

	context := Context{
		ServiceName: name,
		Workspace:   wspConext,

	}
	_, err := config.ReadWorkspaceConfig()
	if err != nil {
		return err
	}

	if err := RenderTemplateTree(context); err != nil {
		return err
	}

	// _, err := workspace.ReadWorkspaceConfig()
	// if err != nil {
	// 	return err
	// }

	// if err := createServiceYaml(name); err != nil {
	// 	return err
	// }

	return nil
}

func createServiceHier(dir string) error {
	fmt.Printf("creating hierarchy in %s\n", dir)

	err := os.Mkdir("backend/cmd/"+dir, 0755)
	if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("failed to create service directory: please remove file or directory with the same name")
	}
	if err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// TODO: README.md
	basePrefixes := []string{
		"schemas",
		"backend/internal",
	}
	for _, prefix := range basePrefixes {
		err = os.MkdirAll(fmt.Sprintf("%s/%s", prefix, dir), 0755)
		if err != nil {
			return fmt.Errorf("failed to create %s/%s directory: %w", prefix, dir, err)
		}
	}

	return nil
}

func createServiceFiles(dir string) error {
	fmt.Printf("creating files in %s\n", dir)

	mainGoRendered := fmt.Sprintf(mainGoTemplate, dir)
	err := ioutil.WriteFile(fmt.Sprintf("backend/cmd/%s/main.go", dir), []byte(mainGoRendered), 0644)
	if err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}
	err = os.MkdirAll(fmt.Sprintf("schemas/%s/api", dir), 0755)
	if err != nil {
		return fmt.Errorf("failed to create schemas dir: %w", err)
	}

	apiRendered := fmt.Sprintf(apiTemplate, dir, "http://"+dir+".company.com")
	err = ioutil.WriteFile(fmt.Sprintf("schemas/%s/api/api.yaml", dir), []byte(apiRendered), 0644)
	if err != nil {
		return fmt.Errorf("failed to create sample api.yaml: %w", err)
	}

	return nil
}

func createServiceYaml(dir string) error {
	fmt.Printf("creating yaml in %s\n", dir)

	conf := config.ServiceConfig{
		ServiceName: dir,
		Maintainers: []string{
			"First maintainer name",
			"Second maintainer name",
		},
	}
	path := fmt.Sprintf("backend/cmd/%s", dir)
	return config.SaveServiceConfig(path, conf)
}

package workspace

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/chebykinn/mify/pkg/mifyconfig"
)

type GoService struct {
	Name string
}

type Context struct {
	Name       string
	BasePath   string
	GoRoot     string // Path to go_services
	Config     mifyconfig.WorkspaceConfig
	TplHeader  string
	GoServices []GoService
}

func InitContext(workspacePath string) (Context, error) {
	if len(workspacePath) == 0 {
		var err error
		workspacePath, err = mifyconfig.FindWorkspaceConfigPath()
		if err != nil {
			return Context{}, err
		}
	}
	conf, err := mifyconfig.ReadWorkspaceConfig(workspacePath)
	if err != nil {
		return Context{}, err
	}

	res := Context{
		Name:      filepath.Base(workspacePath), // TODO: validate
		BasePath:  workspacePath,
		GoRoot:    filepath.Join(workspacePath, mifyconfig.GoServicesRoot),
		Config:    conf,
		TplHeader: "// THIS FILE IS AUTOGENERATED, DO NOT EDIT\n// Generated by mify",
	}

	if err = fillGoServices(&res); err != nil {
		return res, err
	}

	return res, nil
}

// Path to include app.go
func (c Context) GetAppIncludePath(serviceName string) string {
	return fmt.Sprintf(
		"%s/go_services/internal/%s/generated/app",
		c.GetRepository(),
		serviceName)
}

func (c *Context) GetRepository() string {
	return fmt.Sprintf("%s/%s/%s",
		c.Config.GitHost,
		c.Config.GitNamespace,
		c.Config.GitRepository)
}

func fillGoServices(ctx *Context) error {
	ctx.GoServices = make([]GoService, 0)

	files, err := ioutil.ReadDir(filepath.Join(ctx.GoRoot, "cmd"))
	if err != nil {
		return nil
	}

	for _, f := range files {
		if f.IsDir() && f.Name() != "dev-runner" {
			ctx.GoServices = append(ctx.GoServices, GoService{
				Name: f.Name(),
			})
		}
	}

	return nil
}

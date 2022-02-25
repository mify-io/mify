package mify

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/mify-io/mify/internal/mify/userinput"
	"github.com/mify-io/mify/pkg/workspace"
	"github.com/mify-io/mify/pkg/workspace/mutators"
	"github.com/spf13/viper"
)

type Config struct {
	APIToken string `mapstructure:"MIFY_API_TOKEN"`
}

func GetConfigDirectory() string {
	return filepath.Join(xdg.ConfigHome, "mify")
}

func NewDefaultConfig() Config {
	viper.SetDefault("MIFY_API_TOKEN", "")
	return Config{}
}

func SaveConfig(config Config) error {
	viper.Set("MIFY_API_TOKEN", config.APIToken)

	err := os.MkdirAll(GetConfigDirectory(), 0755)
	if err != nil {
		return err
	}
	configPath := filepath.Join(GetConfigDirectory(), "config.yaml")
	return viper.WriteConfigAs(configPath)
}

type CliContext struct {
	Logger        *log.Logger
	Ctx           context.Context
	Cancel        context.CancelFunc
	Config        Config
	WorkspacePath string
	IsVerbose     bool
	UserInput     userinput.UserInput

	workspaceDescription *workspace.Description
	mutatorContext       *mutators.MutatorContext
}

func NewContext(config Config, workspacePath string, isVerbose bool) *CliContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &CliContext{
		Logger:        log.New(os.Stdout, "", 0),
		Ctx:           ctx,
		Cancel:        cancel,
		Config:        config,
		WorkspacePath: workspacePath,
		IsVerbose:     isVerbose,
	}
}

func (c CliContext) GetCtx() context.Context {
	return c.Ctx
}

func (c *CliContext) LoadWorkspace() error {
	res, err := workspace.InitDescription(c.WorkspacePath)
	if err != nil {
		return err
	}

	c.workspaceDescription = &res
	c.WorkspacePath = c.workspaceDescription.BasePath

	c.mutatorContext = mutators.NewMutatorContext(c.Ctx, c.Logger, c.workspaceDescription)

	return nil
}

func (c *CliContext) MustGetWorkspaceDescription() *workspace.Description {
	if c.workspaceDescription == nil {
		panic("missed workspaceDescription")
	}

	return c.workspaceDescription
}

func (c *CliContext) MustGetMutatorContext() *mutators.MutatorContext {
	if c.mutatorContext == nil {
		panic("missed mutatorContext")
	}

	return c.mutatorContext
}

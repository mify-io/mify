package mify

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
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
	Logger *log.Logger
	Ctx    context.Context
	Cancel context.CancelFunc
	Config Config
}

func NewContext(config Config) *CliContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &CliContext{
		Logger: log.New(os.Stdout, "", 0),
		Ctx:    ctx,
		Cancel: cancel,
		Config: config,
	}
}

func initMutatorCtx(ctx *CliContext, basePath string) (*mutators.MutatorContext, error) {
	descr, err := workspace.InitDescription(basePath)
	if err != nil {
		return nil, err
	}

	return mutators.NewMutatorContext(ctx.Ctx, ctx.Logger, &descr), nil
}

func (c CliContext) GetCtx() context.Context {
	return c.Ctx
}

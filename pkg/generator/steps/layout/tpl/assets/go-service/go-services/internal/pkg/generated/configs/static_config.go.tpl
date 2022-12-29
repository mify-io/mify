{{- .Workspace.TplHeader}}
// vim: set ft=go:

package configs

import (
	"reflect"

	"github.com/kelseyhightower/envconfig"
)

type staticConfigLoader struct {
	UseEnv   bool
	UseFiles bool
}

func (l staticConfigLoader) load(cfgType reflect.Type) interface{} {
	newDataPtr := reflect.New(cfgType).Interface()

	if l.UseFiles {
		panic("Not supported yet") // TODO:
	}

	if l.UseEnv {
		// TODO: handle required field. Maybe set required env varibales before call envconfig.Process
		envconfig.MustProcess("", newDataPtr)
	}

	return newDataPtr
}

type MifyStaticConfig struct {
	configProviderBase
}

func NewMifyStaticConfig() (*MifyStaticConfig, error) {
	return &MifyStaticConfig{
		configProviderBase: configProviderBase{
			configs: make(map[string]*storedConfig),
			defaulLoader: staticConfigLoader{
				UseEnv:   true,
				UseFiles: false,
			},
		},
	}, nil
}

type StaticRegisterOpts struct {
	UseFiles bool
	UseEnv   bool
}

func (c *MifyStaticConfig) Register(cfg interface{}, opts StaticRegisterOpts) {
	loader := staticConfigLoader{
		UseEnv:   opts.UseEnv,
		UseFiles: opts.UseFiles,
	}
	c.registerConfig(getConfigType(cfg), loader)
}


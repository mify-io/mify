{{- .Workspace.TplHeader}}

// TODO: support defaults
// TODO: support yaml tags from struct
// TODO: init all configs when start

package configs

import (
	"github.com/hashicorp/consul/api"
	"github.com/lalamove/konfig"
	"github.com/lalamove/konfig/loader/klconsul"
	"github.com/lalamove/konfig/parser/kpyaml"
	"reflect"
)

const (
	configsPath = "config"
)

// Loader

type dynamicConfigLoader struct {
	useConsul    bool
	useFiles     bool
	consulClient *api.Client
	rootStore    konfig.Store
	onUpdate     func(configType reflect.Type, newDataPtr interface{})
}

func (l dynamicConfigLoader) load(cfgType reflect.Type) interface{} {
	store := l.rootStore.Group(cfgType.Name())

	var dataPtr interface{}

	if l.useConsul {
		consulLoader := klconsul.New(&klconsul.Config{
			Client: l.consulClient,
			Keys: []klconsul.Key{
				{
					Key:    configsPath + "/" + cfgType.Name(), // TODO: fuzzy names
					Parser: kpyaml.Parser,
				},
			},
			Watch: true,
		})
		store.RegisterLoader(consulLoader)

		emptyStruct := reflect.Zero(cfgType).Interface()
		store.Bind(emptyStruct)

		err := store.Load()
		if err != nil {
			panic(err)
		}

		newData := reflect.New(cfgType)
		newData.Elem().Set(reflect.ValueOf(store.Value()))
		dataPtr = newData.Interface()

		store.RegisterLoaderWatcher(consulLoader, func(s konfig.Store) error {
			newData := reflect.New(cfgType)
			newData.Elem().Set(reflect.ValueOf(store.Value()))
			l.onUpdate(cfgType, newData.Interface())

			return nil
		})

		store.Watch()
	}

	return dataPtr
}

// MifyDynamicConfig

type MifyDynamicConfig struct {
	configProviderBase
	consulClient *api.Client
}

func NewMifyDynamicConfig(consulClient *api.Client) (*MifyDynamicConfig, error) {
	loader := dynamicConfigLoader{
		useConsul:    true,
		useFiles:     false,
		consulClient: consulClient,
		rootStore:    konfig.New(konfig.DefaultConfig()),
	}
	dynamicConfig := &MifyDynamicConfig{
		configProviderBase: configProviderBase{
			configs:      make(map[string]*storedConfig),
			defaulLoader: loader,
		},
	}
	loader.onUpdate = dynamicConfig.updateConfigData

	return dynamicConfig, nil
}

type DynamicRegisterOpts struct {
	UseConsul bool
	UseFiles  bool
}

func (c *MifyDynamicConfig) RegisterConfig(cfgType interface{}, opts DynamicRegisterOpts) {
	loader := dynamicConfigLoader{
		useConsul:    opts.UseConsul,
		useFiles:     opts.UseFiles,
		consulClient: c.consulClient,
		rootStore:    konfig.New(konfig.DefaultConfig()),
		onUpdate:     c.updateConfigData,
	}
	c.registerConfig(getConfigType(cfgType), loader)
}


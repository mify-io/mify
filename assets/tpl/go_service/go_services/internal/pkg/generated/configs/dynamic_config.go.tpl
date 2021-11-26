{{- .Workspace.TplHeader}}

// TODO: support defaults
// TODO: support yaml tags from struct
// TODO: init all configs when start

package configs

import (
	"reflect"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/lalamove/konfig"
	"github.com/lalamove/konfig/loader/klconsul"
	"github.com/lalamove/konfig/parser/kpyaml"
)

type MyConfig struct {
	Data string `yaml:"data" default:"127.0.0.1:8500"`
}

const (
	configsPath = "config"
)

var defaultDynamicRegisterOpts = DynamicRegisterOpts{
	UseConsul: true,
	UseFiles:  false,
}

type registeredDynamicConfig struct {
	structType reflect.Type
	data       interface{}
	opts       DynamicRegisterOpts
	store      konfig.Store
}

type MifyDynamicConfig struct {
	configs        map[string]*registeredDynamicConfig // TypeName -> config
	configsRwMutex sync.RWMutex

	rootStore    konfig.Store
	consulClient *api.Client
}

func NewMifyDynamicConfig(consulClient *api.Client) (*MifyDynamicConfig, error) {
	// TODO: start watching all configs in config/
	return &MifyDynamicConfig{
		configs:      make(map[string]*registeredDynamicConfig),
		rootStore:    konfig.New(konfig.DefaultConfig()),
		consulClient: consulClient,
	}, nil
}

type DynamicRegisterOpts struct {
	UseConsul bool
	UseFiles  bool
}

func (c *MifyDynamicConfig) addConfigImpl(cfgType reflect.Type, opts DynamicRegisterOpts) *registeredDynamicConfig {
	cfg := &registeredDynamicConfig{
		structType: cfgType,
		data:       reflect.New(cfgType).Interface(),
		store:      c.rootStore.Group(cfgType.Name()),
		opts:       opts,
	}

	if cfg.opts.UseConsul {
		consulLoader := klconsul.New(&klconsul.Config{
			Client: c.consulClient,
			Keys: []klconsul.Key{
				{
					Key:    configsPath + "/" + cfgType.Name(), // TODO: fuzzy names
					Parser: kpyaml.Parser,
				},
			},
			Watch: true,
		})
		cfg.store.RegisterLoader(consulLoader)

		emptyStruct := reflect.Zero(cfg.structType).Interface()
		cfg.store.Bind(emptyStruct)

		err := cfg.store.Load()
		if err != nil {
			panic(err)
		}

		newData := reflect.New(cfg.structType)
		newData.Elem().Set(reflect.ValueOf(cfg.store.Value()))
		cfg.data = newData.Interface()

		cfg.store.RegisterLoaderWatcher(consulLoader, func(s konfig.Store) error {
			c.configsRwMutex.Lock()
			defer func() { c.configsRwMutex.Unlock() }()

			newData := reflect.New(cfg.structType)
			newData.Elem().Set(reflect.ValueOf(cfg.store.Value()))
			cfg.data = newData.Interface()

			return nil
		})

		cfg.store.Watch()
	}

	c.configs[cfgType.Name()] = cfg

	if cfg.opts.UseFiles {
		panic("Not supported yet") // TODO:
	}

	return cfg
}

func (c *MifyDynamicConfig) safeGetConfig(cfgType reflect.Type) *registeredDynamicConfig {
	c.configsRwMutex.RLock()
	defer func() { c.configsRwMutex.RUnlock() }()

	if val, ok := c.configs[cfgType.Name()]; ok {
		return val
	}

	return nil
}

func (c *MifyDynamicConfig) addOrGetConfig(cfgType reflect.Type, opts DynamicRegisterOpts) *registeredDynamicConfig {
	cfg := c.safeGetConfig(cfgType)
	if cfg != nil {
		return cfg
	}

	c.configsRwMutex.Lock()
	defer func() { c.configsRwMutex.Unlock() }()

	if val, ok := c.configs[cfgType.Name()]; ok {
		return val
	}

	return c.addConfigImpl(cfgType, opts)
}

// Public

func (c *MifyDynamicConfig) Register(cfgPtr interface{}, opts DynamicRegisterOpts) {
	configType := getConfigType(cfgPtr)

	// TODO: rewrite existing?
	c.addOrGetConfig(configType, opts)
}

// Returns ptr to config
func (c *MifyDynamicConfig) Get(cfgPtr interface{}) (interface{}, error) {
	configType := getConfigType(cfgPtr)

	cfg := c.addOrGetConfig(configType, defaultDynamicRegisterOpts)

	return cfg.data, nil
}

func (c *MifyDynamicConfig) MustGet(cfgPtr interface{}) interface{} {
	res, err := c.Get(cfgPtr)
	if err != nil {
		panic(err)
	}
	return res
}

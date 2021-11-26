{{- .Workspace.TplHeader}}

package configs

import (
	"reflect"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type registeredConfig struct {
	structType reflect.Type
	data       interface{}
	opts       RegisterOpts
}

type MifyStaticConfig struct {
	configs        map[string]*registeredConfig // TypeName -> config
	configsRwMutex sync.RWMutex
}

func NewMifyStaticConfig() (*MifyStaticConfig, error) {
	return &MifyStaticConfig{
		configs: make(map[string]*registeredConfig),
	}, nil
}

type RegisterOpts struct {
	UseFiles bool
	UseEnv   bool
}

func bootstrapConfig(cfg *registeredConfig) {
	cfg.data = reflect.New(cfg.structType).Interface()

	if cfg.opts.UseFiles {
		panic("Not supported yet") // TODO:
	}

	if cfg.opts.UseEnv {
		// TODO: handle required field. Maybe set required env varibales before call envconfig.Process
		envconfig.MustProcess("", cfg.data)
	}
}

func (c *MifyStaticConfig) addConfigImpl(cfgType reflect.Type, opts RegisterOpts) *registeredConfig {
	cfg := &registeredConfig{
		structType: cfgType,
		data:       nil,
		opts:       opts,
	}
	c.configs[cfgType.Name()] = cfg

	bootstrapConfig(cfg)
	return cfg
}

func (c *MifyStaticConfig) safeGetConfig(cfgType reflect.Type) *registeredConfig {
	c.configsRwMutex.RLock()
	defer func() { c.configsRwMutex.RUnlock() }()

	if val, ok := c.configs[cfgType.Name()]; ok {
		return val
	}

	return nil
}

func (c *MifyStaticConfig) addOrGetConfig(cfgType reflect.Type, opts RegisterOpts) *registeredConfig {
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

func (c *MifyStaticConfig) Register(cfgPtr interface{}, opts RegisterOpts) {
	configType := getConfigType(cfgPtr)

	// TODO: rewrite existing?
	c.addOrGetConfig(configType, opts)
}

// Returns ptr to config
func (c *MifyStaticConfig) Get(cfgPtr interface{}) (interface{}, error) {
	configType := getConfigType(cfgPtr)

	cfg := c.addOrGetConfig(configType, RegisterOpts{
		UseFiles: false,
		UseEnv:   true,
	})

	return cfg.data, nil
}

func (c *MifyStaticConfig) MustGet(cfgPtr interface{}) interface{} {
	res, err := c.Get(cfgPtr)
	if err != nil {
		panic(err)
	}
	return res
}

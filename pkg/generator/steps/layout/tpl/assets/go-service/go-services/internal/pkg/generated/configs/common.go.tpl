{{- .Workspace.TplHeader}}
// vim: set ft=go:

package configs

import (
	"reflect"
	"sync"
)

func getConfigTypeByPtr(cfgPtr interface{}) reflect.Type {
	configPtrType := reflect.TypeOf(cfgPtr)
	if configPtrType.Kind() != reflect.Ptr || configPtrType.Elem().Kind() != reflect.Struct {
		panic("must be a pointer to struct")
	}
	return configPtrType.Elem()
}

func getConfigType(cfg interface{}) reflect.Type {
	configType := reflect.TypeOf(cfg)
	if configType.Kind() != reflect.Struct {
		panic("must be a struct")
	}
	return configType
}

type cfgLoader interface {
	load(cfgType reflect.Type) interface{}
}

// Store

type storedConfig struct {
	dataPtr interface{}
	loader  cfgLoader
}

type configProviderBase struct {
	configs        map[string]*storedConfig // TypeName -> config
	configsRwMutex sync.RWMutex
	defaulLoader   cfgLoader
}

func (c *configProviderBase) doUnderReadLock(f func()) {
	c.configsRwMutex.RLock()
	defer c.configsRwMutex.RUnlock()
	f()
}

func (c *configProviderBase) doUnderWriteLock(f func()) {
	c.configsRwMutex.Lock()
	defer func() { c.configsRwMutex.Unlock() }()
	f()
}

func (c *configProviderBase) registerConfig(cfgType reflect.Type, loader cfgLoader) {
	c.doUnderWriteLock(func() {
		c.configs[cfgType.Name()] = &storedConfig{
			dataPtr: nil,
			loader:  loader,
		}
	})
}

func (c *configProviderBase) getConfigDataPtr(cfgType reflect.Type) interface{} {
	if val, ok := c.configs[cfgType.Name()]; ok {
		return val.dataPtr
	}

	return nil
}

func (c *configProviderBase) addOrGetConfig(configType reflect.Type) interface{} {
	var cfgDataPtr interface{}

	c.doUnderReadLock(func() { cfgDataPtr = c.getConfigDataPtr(configType) })
	if cfgDataPtr != nil {
		return cfgDataPtr
	}

	c.doUnderWriteLock(func() {
		if val, ok := c.configs[configType.Name()]; ok {
			if val.dataPtr != nil {
				cfgDataPtr = val.dataPtr
				return
			}

			val.dataPtr = val.loader.load(configType)
			cfgDataPtr = val.dataPtr
			return
		}

		storedConfigPtr := &storedConfig{
			dataPtr: c.defaulLoader.load(configType),
			loader:  c.defaulLoader,
		}
		c.configs[configType.Name()] = storedConfigPtr
		cfgDataPtr = storedConfigPtr.dataPtr
	})

	return cfgDataPtr
}

func (c *configProviderBase) updateConfigData(configType reflect.Type, newDataPtr interface{}) {
	c.doUnderWriteLock(func() {
		if val, ok := c.configs[configType.Name()]; ok {
			val.dataPtr = newDataPtr
			return
		}

		c.configs[configType.Name()] = &storedConfig{
			dataPtr: newDataPtr,
			loader:  c.defaulLoader,
		}
	})
}

// Public

// Extracts a type of struct from cfgPtr and returns a ptr to stored config for this type.
// cfgPtr must be reference to configuration struct (struct wich will contains configuration).
// This function has better performance than func Get, because no excess struct instances are created.
// Example:
//   type MyStruct struct { SomeData string }
//
//   var c configProviderBase
//   cfgPtr := c.GetPtr((*MyStruct)(nil)).(*MyStruct)
func (c *configProviderBase) GetPtr(cfgPtr interface{}) (interface{}, error) {
	configType := getConfigTypeByPtr(cfgPtr)
	cfgDataPtr := c.addOrGetConfig(configType)
	return cfgDataPtr, nil
}

// Same as GetPtr, but panics when any error is occurred
func (c *configProviderBase) MustGetPtr(cfgPtr interface{}) interface{} {
	res, err := c.GetPtr(cfgPtr)
	if err != nil {
		panic(err)
	}
	return res
}

// Extracts a type of struct from cfg and returns a stored config for this type.
// cfg must be configuration struct (struct wich will contains configuration).
// This function has worse performance than func GetPtr, because excess struct instances are created.
// Example:
//   type MyStruct struct { SomeData string }
//
//   var c configProviderBase
//   cfg := c.Get(MyStruct{}).(MyStruct)
func (c *configProviderBase) Get(cfg interface{}) (interface{}, error) {
	configType := getConfigType(cfg)
	cfgDataPtr := c.addOrGetConfig(configType)
	return reflect.ValueOf(cfgDataPtr).Elem().Interface(), nil
}

// Same as Get, but panics when any error is occurred
func (c *configProviderBase) MustGet(cfg interface{}) interface{} {
	res, err := c.Get(cfg)
	if err != nil {
		panic(err)
	}
	return res
}

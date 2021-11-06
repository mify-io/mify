package core

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bketelsen/crypt/backend/consul"
	"github.com/bketelsen/crypt/config"
	"go.uber.org/zap"
)

const (
	configsPath   = "config"
	refreshPeriod = 60 * time.Second
)

type MifyDynamicConfig struct {
	rwMutex sync.RWMutex
	data    map[string][]byte

	manager            config.ConfigManager
	stopChannels       map[string]chan bool
	mifyServiceContext *MifyServiceContext
	logger             *zap.Logger
}

func NewMifyDynamicConfig(mifyServiceContext *MifyServiceContext) (*MifyDynamicConfig, error) {
	cl, err := consul.New([]string{"127.0.0.1:8500"}) // TODO: to config
	if err != nil {
		return nil, err
	}

	cm, err := config.NewStandardConfigManager(cl)
	if err != nil {
		return nil, err
	}

	ctxWrapper := &MifyDynamicConfig{
		data:               make(map[string][]byte),
		manager:            cm,
		stopChannels:       make(map[string]chan bool),
		mifyServiceContext: mifyServiceContext,
		logger:             mifyServiceContext.LoggerFor("dynamic_config"),
	}

	ctxWrapper.fullRefresh()

	ticker := time.NewTicker(refreshPeriod)
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				ctxWrapper.fullRefresh()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	return ctxWrapper, nil
}

func (cw *MifyDynamicConfig) fullRefresh() error {
	cw.logger.Info("Start updating configs...")

	configs, err := cw.manager.List(configsPath)
	if err != nil {
		return err
	}

	cw.logger.Sugar().Infof("Loaded %s configs", len(configs))

	knownKeys := make(map[string]struct{})

	cw.rwMutex.Lock()
	defer func() { cw.rwMutex.Unlock() }()

	for _, config := range configs {
		cfgName := strings.Replace(config.Key, configsPath+"/", "", 1)
		knownKeys[cfgName] = struct{}{}

		_, ok := cw.data[cfgName]
		if ok {
			// Config is already under watch
			continue
		}

		cw.registerConfig(cfgName, config.Value)
	}

	for cfgName := range cw.data {
		if _, ok := knownKeys[cfgName]; !ok {
			cw.unregisterConfig(cfgName)
		}
	}

	return nil
}

// Should be called under write lock
func (cw *MifyDynamicConfig) registerConfig(cfgName string, data []byte) {
	cw.logger.Sugar().Infof("Registering new config %s ...", cfgName)

	cw.data[cfgName] = data

	stop := make(chan bool)
	cw.stopChannels[cfgName] = stop

	resp := cw.manager.Watch(getConfigPath(cfgName), stop)

	go func() {
		for {
			r := <-resp
			if r.Error != nil {
				cw.logger.Sugar().Error(r.Error)
				continue
			}

			cw.logger.Sugar().Infof("Config %s is changed", cfgName)

			cw.rwMutex.Lock()

			if _, ok := cw.data[cfgName]; !ok {
				cw.rwMutex.Unlock()
				// It means that config was unregistered. And we should stop polling
				break
			}

			cw.data[cfgName] = r.Value
			cw.rwMutex.Unlock()
		}
	}()
}

// Should be called under write lock
func (cw *MifyDynamicConfig) unregisterConfig(cfgName string) {
	cw.logger.Sugar().Infof("Unregistering config %s ...", cfgName)

	delete(cw.data, cfgName)
	delete(cw.stopChannels, cfgName)

	// Bug inside crypt library. Channel is not using at all
	//cw.stopChannels[name] <- true
}

func getConfigPath(cfgName string) string {
	return configsPath + "/" + cfgName
}

func (cw *MifyDynamicConfig) GetConfig(cfgName string) ([]byte, error) {
	cw.rwMutex.RLock()
	defer func() { cw.rwMutex.RUnlock() }()

	cfg, ok := cw.data[cfgName]
	if !ok {
		return nil, fmt.Errorf("сonfig with name %s wasn't found", cfgName)
	}

	return cfg, nil
}

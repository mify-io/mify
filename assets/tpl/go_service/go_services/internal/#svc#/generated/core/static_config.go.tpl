package core

import (
	"github.com/lalamove/konfig"
	"github.com/lalamove/konfig/loader/klenv"
)

type MifyStaticConfig struct {
	store konfig.Store
}

func NewMifyStaticConfig() (*MifyStaticConfig, error) {
	store := konfig.New(konfig.DefaultConfig())
	store.RegisterLoader(klenv.New(&klenv.Config{}))
	store.Load()
	return &MifyStaticConfig{
		store: store,
	}, nil
}

func (c *MifyStaticConfig) Get(key string) interface{} {
	return c.store.Get(key)
}

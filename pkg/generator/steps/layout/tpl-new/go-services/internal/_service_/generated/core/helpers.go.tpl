{{- .TplHeader}}
// vim: set ft=go:

package core

import "{{.ConfigsImportPath}}"

type StaticConfigProvider interface {
	StaticConfig() *configs.MifyStaticConfig
}

var _ StaticConfigProvider = (*MifyServiceContext)(nil)
var _ StaticConfigProvider = (*MifyRequestContext)(nil)

func GetStaticConfig[TConfig any](ctx StaticConfigProvider) (*TConfig, error) {
	res, err := ctx.StaticConfig().GetPtr((*TConfig)(nil))
	if err != nil {
		return nil, err
	}

	return res.(*TConfig), nil
}

func MustGetStaticConfig[TConfig any](ctx StaticConfigProvider) *TConfig {
	res, err := GetStaticConfig[TConfig](ctx)
	if err != nil {
		panic(err)
	}
	return res
}

type DynamicConfigProvider interface {
	DynamicConfig() *configs.MifyDynamicConfig
}

var _ DynamicConfigProvider = (*MifyServiceContext)(nil)
var _ DynamicConfigProvider = (*MifyRequestContext)(nil)

func GetDynamicConfig[TConfig any](ctx DynamicConfigProvider) (*TConfig, error) {
	res, err := ctx.DynamicConfig().GetPtr((*TConfig)(nil))
	if err != nil {
		return nil, err
	}

	return res.(*TConfig), nil
}

func MustGetDynamicConfig[TConfig any](ctx DynamicConfigProvider) *TConfig {
	res, err := GetDynamicConfig[TConfig](ctx)
	if err != nil {
		panic(err)
	}
	return res
}
{{- .TplHeader}}

package core

func GetStaticConfig[TConfig any](ctx *MifyRequestContext) (*TConfig, error) {
	res, err := ctx.StaticConfig().GetPtr((*TConfig)(nil))
	if err != nil {
		return nil, err
	}

	return res.(*TConfig), nil
}

func MustGetStaticConfig[TConfig any](ctx *MifyRequestContext) *TConfig {
	res, err := GetStaticConfig[TConfig](ctx)
	if err != nil {
		panic(err)
	}
	return res
}

func GetDynamicConfig[TConfig any](ctx *MifyRequestContext) (*TConfig, error) {
	res, err := ctx.DynamicConfig().GetPtr((*TConfig)(nil))
	if err != nil {
		return nil, err
	}

	return res.(*TConfig), nil
}

func MustGetDynamicConfig[TConfig any](ctx *MifyRequestContext) *TConfig {
	res, err := GetDynamicConfig[TConfig](ctx)
	if err != nil {
		panic(err)
	}
	return res
}
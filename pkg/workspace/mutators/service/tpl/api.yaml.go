package tpl

type ApiSchemaModel struct {
	ServiceName string
}

func NewApiSchemaModel(serviceName string) ApiSchemaModel {
	return ApiSchemaModel{
		ServiceName: serviceName,
	}
}

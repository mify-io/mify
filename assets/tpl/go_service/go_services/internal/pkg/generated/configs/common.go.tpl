{{- .Workspace.TplHeader}}

package configs

import "reflect"

func getConfigType(cfgPtr interface{}) reflect.Type {
	configPtrType := reflect.TypeOf(cfgPtr)
	if configPtrType.Kind() != reflect.Ptr {
		panic("must be a pointer")
	}
	return configPtrType.Elem()
}

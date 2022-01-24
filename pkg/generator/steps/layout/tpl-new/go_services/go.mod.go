package goservices

type GoModModel struct {
	ModuleName string
}

func NewGoModModel(moduleName string) GoModModel {
	return GoModModel{
		ModuleName: moduleName,
	}
}

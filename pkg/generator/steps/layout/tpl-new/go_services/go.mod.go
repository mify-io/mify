package goservices

type goModModel struct {
	ModuleName string
}

func newGoModModel(moduleName string) goModModel {
	return goModModel{
		ModuleName: moduleName,
	}
}

package tpl

type TargetService struct {
	Name           string
	SafeName       string
	AppIncludePath string
}

func NewTargetService(name string, safeName string, appIncludePath string) TargetService {
	return TargetService{
		Name:           name,
		SafeName:       safeName,
		AppIncludePath: appIncludePath,
	}
}

type Model struct {
	Header   string
	Services []TargetService
}

func NewModel(header string, services []TargetService) Model {
	return Model{
		Header:   header,
		Services: services,
	}
}

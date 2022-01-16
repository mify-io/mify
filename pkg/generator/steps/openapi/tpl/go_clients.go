package tpl

type GoClientModel struct {
	ClientName       string
	PackageName      string
	PrivateFieldName string
	PublicMethodName string
	IncludePath      string
}

func NewGoClientModel(
	clientName string,
	packageName string,
	privateFieldName string,
	publicMethodName string,
	includePath string) GoClientModel {

	return GoClientModel{
		ClientName:       clientName,
		PackageName:      packageName,
		PrivateFieldName: privateFieldName,
		PublicMethodName: publicMethodName,
		IncludePath:      includePath,
	}
}

type GoClientsModel struct {
	Header      string
	GoModule    string
	ServiceName string

	MetricsIncludePath string
	Clients            []GoClientModel
}

func NewGoClientsModel(
	header string,
	goModule string,
	serviceName string,
	metricsIncludePath string,
	clients []GoClientModel) GoClientsModel {

	return GoClientsModel{
		Header:      header,
		GoModule:    goModule,
		ServiceName: serviceName,

		MetricsIncludePath: metricsIncludePath,
		Clients:            clients,
	}
}

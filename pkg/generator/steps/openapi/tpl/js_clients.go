package tpl

type JsClientModel struct {
	ClientName       string
	PublicMethodName string
	ClassName        string
}

func NewJsClientModel(
	clientName string,
	className string,
	publicMethodName string) JsClientModel {

	return JsClientModel{
		ClientName:       clientName,
		PublicMethodName: publicMethodName,
		ClassName:        className,
	}
}

type JsClientsModel struct {
	Header      string
	ServiceName string
	Clients     []JsClientModel
}

func NewJsClientsModel(
	header string,
	serviceName string,
	clients []JsClientModel) JsClientsModel {

	return JsClientsModel{
		Header:      header,
		ServiceName: serviceName,
		Clients:     clients,
	}
}

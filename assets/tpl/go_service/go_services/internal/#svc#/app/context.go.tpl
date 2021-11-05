{{- .Workspace.TplHeader}}

package app

type ServiceContext struct {
	// Append your dependencies here
}

func NewServiceContext () (*ServiceContext, error) {
	context := &ServiceContext{
		// Here you can initialize your dependencies
	}
	return context, nil
}

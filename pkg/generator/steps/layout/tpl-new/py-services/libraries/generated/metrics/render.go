package metrics

import (
	_ "embed"
	"path/filepath"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed server.py.tpl
var serverTemplate string

//go:embed request_metrics.py.tpl
var requestMetricsTemplate string

//go:embed client_metrics.py.tpl
var clientMetricsTemplate string

func Render(ctx *gencontext.GenContext) error {
	serverModel := struct{}{}
	serverPath := filepath.Join(ctx.GetWorkspace().GetPythonServicesLibrariesGeneratedMetricsAbsPath(), "server.py")
	if err := render.RenderOrSkipTemplate(serverTemplate, serverModel, serverPath); err != nil {
		return render.WrapError("server.py", err)
	}

	requestMetricsModel := newRequestMetricsModel(ctx)
	requestMetricsPath := filepath.Join(ctx.GetWorkspace().GetPythonServicesLibrariesGeneratedMetricsAbsPath(), "request_metrics.py")
	if err := render.RenderOrSkipTemplate(requestMetricsTemplate, requestMetricsModel, requestMetricsPath); err != nil {
		return render.WrapError("request_metrics.py", err)
	}

	clientMetricsModel := newClientMetricsModel(ctx)
	clientMetricsPath := filepath.Join(ctx.GetWorkspace().GetPythonServicesLibrariesGeneratedMetricsAbsPath(), "client_metrics.py")
	if err := render.RenderOrSkipTemplate(clientMetricsTemplate, clientMetricsModel, clientMetricsPath); err != nil {
		return render.WrapError("client_metrics.py", err)
	}

	return nil
}

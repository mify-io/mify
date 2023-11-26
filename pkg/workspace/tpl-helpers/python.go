package tplhelpers

import (
	"path"
	"strings"

	"github.com/mify-io/mify/pkg/mifyconfig"
)

type pythonHelpers struct {}

func (h pythonHelpers) MakeDefaultMifyGeneratedPath() string {
	return path.Join(mifyconfig.PythonServicesRoot, mifyconfig.MifyGeneratedDirName)
}

func (h pythonHelpers) MakeDefaultMifyGeneratedPackage(conf mifyconfig.WorkspaceConfig, generatedPath string) string {
	return strings.ReplaceAll(mifyconfig.MifyGeneratedDirName, "-", "_")
}

package tplhelpers

import (
	"fmt"
	"path"
	"strings"

	"github.com/mify-io/mify/pkg/mifyconfig"
)

type pythonHelpers struct {}

func (h pythonHelpers) MakeDefaultMifyGeneratedPath() string {
	dirName := strings.ReplaceAll(mifyconfig.MifyGeneratedDirName, "-", "_")
	return path.Join(mifyconfig.PythonServicesRoot, dirName)
}

func (h pythonHelpers) MakeDefaultMifyGeneratedPackage(conf mifyconfig.WorkspaceConfig, generatedPath string) string {
	pkgName := strings.ReplaceAll(mifyconfig.MifyGeneratedDirName, "/", ".")
	pkgName = strings.ReplaceAll(pkgName, "-", "_")
	return pkgName
}

func (h pythonHelpers) GetCommonPackage(pkgRoot string) string {
	return fmt.Sprintf("%s.common", pkgRoot)
}

func (h pythonHelpers) GetServicePackage(pkgRoot string, serviceName string) string {
	return fmt.Sprintf("%s.services.%s", pkgRoot, serviceName)
}

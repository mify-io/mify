package tplhelpers

import (
	"fmt"
	"path"

	"github.com/mify-io/mify/pkg/mifyconfig"
)

type goHelpers struct {}

func (h goHelpers) MakeDefaultMifyGeneratedPath() string {
	return path.Join(mifyconfig.GoServicesRoot, "internal", mifyconfig.MifyGeneratedDirName)
}

func getRepository(conf mifyconfig.WorkspaceConfig) string {
	return fmt.Sprintf("%s/%s/%s",
		conf.GitHost,
		conf.GitNamespace,
		conf.GitRepository)
}

func (h goHelpers) MakeDefaultMifyGeneratedPackage(conf mifyconfig.WorkspaceConfig, generatedPath string) string {
	return fmt.Sprintf("%s/%s", getRepository(conf), generatedPath)
}

func (h goHelpers) GetCommonPackage(pkgRoot string) string {
	return fmt.Sprintf("%s/common", pkgRoot)
}

func (h goHelpers) GetServicePackage(pkgRoot string, serviceName string) string {
	return fmt.Sprintf("%s/services/%s", pkgRoot, serviceName)
}

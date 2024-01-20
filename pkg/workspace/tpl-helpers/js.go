package tplhelpers

import (
	"fmt"
	"path"

	"github.com/mify-io/mify/pkg/mifyconfig"
)

type jsHelpers struct { }

func (h jsHelpers) MakeDefaultMifyGeneratedPath() string {
	return path.Join(mifyconfig.JsServicesRoot, mifyconfig.MifyGeneratedDirName)
}

func (h jsHelpers) MakeDefaultMifyGeneratedPackage(conf mifyconfig.WorkspaceConfig, generatedPath string) string {
	return mifyconfig.MifyGeneratedDirName
}

func (h jsHelpers) GetCommonPackage(pkgRoot string) string {
	return fmt.Sprintf("%s.common", pkgRoot)
}

func (h jsHelpers) GetServicePackage(pkgRoot string, serviceName string) string {
	return fmt.Sprintf("%s.services.%s", pkgRoot, serviceName)
}

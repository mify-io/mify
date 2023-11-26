package tplhelpers

import (
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

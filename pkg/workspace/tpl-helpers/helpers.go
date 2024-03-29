package tplhelpers

import (
	"github.com/mify-io/mify/pkg/mifyconfig"
)

type WorkspaceTemplateHelpers interface {
	MakeDefaultMifyGeneratedPath() string
	MakeDefaultMifyGeneratedPackage(c mifyconfig.WorkspaceConfig, generatedPath string) string
	GetCommonPackage(pkgRoot string) string
	GetServicePackage(pkgRoot string, serviceName string) string
}

var HelpersMap = map[mifyconfig.ServiceLanguage]WorkspaceTemplateHelpers {
	mifyconfig.ServiceLanguageGo: goHelpers{},
	mifyconfig.ServiceLanguageJs: jsHelpers{},
	mifyconfig.ServiceLanguagePython: pythonHelpers{},
}

func PopulateGeneratorParams(c mifyconfig.WorkspaceConfig) mifyconfig.GeneratorParams {
	params := c.GeneratorParams
	if params.Template == nil {
		params.Template = make(map[mifyconfig.ServiceLanguage]mifyconfig.GeneratorTemplateParams)
	}
	for lang, helper := range HelpersMap {
		tplParams := params.Template[lang]
		if len(tplParams.MifyGeneratedPath) == 0 {
			tplParams.MifyGeneratedPath = helper.MakeDefaultMifyGeneratedPath()
		}
		if len(tplParams.MifyGeneratedPackage) == 0 {
			tplParams.MifyGeneratedPackage = helper.MakeDefaultMifyGeneratedPackage(c, tplParams.MifyGeneratedPath)
		}
		if tplParams.DevRunner == nil {
			tplParams.DevRunner = mifyconfig.MakeDefaultComponent()
		}
		params.Template[lang] = tplParams
	}
	return params
}

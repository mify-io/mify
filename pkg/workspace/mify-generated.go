package workspace

import (
	"fmt"
	"path"

	"github.com/mify-io/mify/pkg/mifyconfig"
)

type MifyGenerated struct {
	lang mifyconfig.ServiceLanguage
	serviceName string
	desc Description
}

func newMifyGenerated(desc Description, serviceConfig *mifyconfig.ServiceConfig) MifyGenerated {
	return MifyGenerated{
		lang: serviceConfig.Language,
		serviceName: serviceConfig.ServiceName,
		desc: desc,
	}
}

func (c MifyGenerated) GetPath() Path {
	return NewPath(c.desc, c.desc.Config.GeneratorParams.Template[c.lang].MifyGeneratedPath)
}

func (c MifyGenerated) GetCommonPath() Path {
	return NewPath(c.desc, path.Join(c.GetPath().Rel(), "common"))
}

func (c MifyGenerated) GetServicePath() Path {
	return NewPath(c.desc, path.Join(c.GetPath().Rel(), "services", c.serviceName))
}

func (c MifyGenerated) GetPackage() string {
	return c.desc.Config.GeneratorParams.Template[c.lang].MifyGeneratedPackage
}

func (c MifyGenerated) GetCommonPackage() string {
	return fmt.Sprintf("%s/common", c.GetPackage())
}

func (c MifyGenerated) GetServicePackage() string {
	return fmt.Sprintf("%s/services/%s", c.GetPackage(), c.serviceName)
}

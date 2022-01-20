package mifyconfig

import (
	"os"
	"path/filepath"
)

type Workspace struct {
	Name       string
	Directory  string
	ConfigPath string
}

func NewWorkspace(workspaceDir string) (Workspace, error) {
	confPath, err := FindWorkspaceConfigPathInLocation(workspaceDir)
	if err != nil {
		return Workspace{}, err
	}
	name := filepath.Base(workspaceDir)
	return Workspace{
		Name:       name,
		Directory:  workspaceDir,
		ConfigPath: confPath,
	}, nil
}

func (w Workspace) GetServicesForLanguage(lang ServiceLanguage) ([]Service, error) {
	svcList := []Service{}
	subDir, _ := getServiceConfigPathByLang(lang)
	globStr := filepath.Join(w.Directory, subDir, "*")
	services, err := filepath.Glob(globStr)
	if err != nil {
		return []Service{}, err
	}
	for _, f := range services {
		stat, err := os.Stat(f)
		if err != nil {
			return []Service{}, err
		}
		if !stat.IsDir() {
			continue
		}
		svc, err := NewService(w, filepath.Base(f))
		if err != nil {
			return []Service{}, err
		}
		svcList = append(svcList, svc)
	}
	return svcList, nil
}

func (w Workspace) GetServices() ([]Service, error) {
	svcList := []Service{}
	for _, l := range LanguagesList {
		langList, err := w.GetServicesForLanguage(l)
		if err != nil {
			return []Service{}, err
		}
		svcList = append(svcList, langList...)
	}
	return svcList, nil
}

func (w Workspace) HasService(serviceName string) (bool, error) {
	svcList, err := w.GetServices()
	if err != nil {
		return false, err
	}
	for _, svc := range svcList {
		if svc.Name == serviceName {
			return true, nil
		}
	}
	return false, nil
}

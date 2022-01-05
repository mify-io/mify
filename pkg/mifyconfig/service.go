package mifyconfig

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	ServiceConfigName = "service.mify.yaml"
)

type ServiceLanguage string

const (
	ServiceLanguageUnknown ServiceLanguage = "unknown"
	ServiceLanguageGo      ServiceLanguage = "go"
	ServiceLanguageJs      ServiceLanguage = "js"
)

var LanguagesList = []ServiceLanguage{
	ServiceLanguageGo,
	ServiceLanguageJs,
}

type Service struct {
	Name string
	Language ServiceLanguage
	Directory string
	ConfigPath string
}


func NewService(workspace Workspace, serviceName string) (Service, error) {
	serviceDir, lang, err := locateServiceDirectory(workspace.Directory, serviceName)
	if err != nil {
		return Service{}, fmt.Errorf("failed to read service %s: %w", serviceName, err)
	}
	configFile := filepath.Join(serviceDir, ServiceConfigName)
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		return Service{}, fmt.Errorf("failed to read service %s: %s: not found", serviceName, configFile)
	}
	return Service{
		Name: filepath.Base(serviceName),
		Language: lang,
		Directory: serviceDir,
		ConfigPath: configFile,
	}, nil
}

func (s Service) ReadConfig() (ServiceConfig, error) {
	rawData, err := ioutil.ReadFile(s.ConfigPath)
	if err != nil {
		return ServiceConfig{}, fmt.Errorf("failed to read service config: %w", err)
	}

	var data ServiceConfig
	err = yaml.Unmarshal(rawData, &data)
	if err != nil {
		return ServiceConfig{}, fmt.Errorf("failed to read service config: %w", err)
	}
	data.Language = s.Language

	return data, nil
}

func getServiceConfigPathByLang(language ServiceLanguage) (string, error) {
	switch language {
	case ServiceLanguageGo:
		return GoServicesRoot + "/cmd", nil
	case ServiceLanguageJs:
		return JsServicesRoot, nil
	}
	return "", fmt.Errorf("no such language: %s", language)
}

func locateServiceDirectory(workspaceDir string, serviceName string) (string, ServiceLanguage, error) {
	for _, l := range LanguagesList {
		subDir, _ := getServiceConfigPathByLang(l)
		dir := filepath.Join(workspaceDir, subDir, serviceName)
		d, err := os.Stat(dir);
		if err == nil && d.IsDir() {
			return dir, l, nil
		}
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		return "", ServiceLanguageUnknown, err
	}
	return "", ServiceLanguageUnknown, fmt.Errorf("service %s not found", serviceName)
}


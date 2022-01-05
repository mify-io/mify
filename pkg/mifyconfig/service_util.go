package mifyconfig

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func getServiceConfigPathByLang(language ServiceLanguage) (string, error) {
	switch language {
	case ServiceLanguageGo:
		return GoServicesRoot + "/cmd", nil
	case ServiceLanguageJs:
		return JsServicesRoot, nil
	}
	return "", fmt.Errorf("no such language: %s", language)
}

func locateServiceConfigPath(workspaceDir string, serviceName string) (string, ServiceLanguage, error) {
	for _, l := range LanguagesList {
		subDir, _ := getServiceConfigPathByLang(l)
		configFile := filepath.Join(workspaceDir, subDir, serviceName, ServiceConfigName)
		if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
			continue
		}
		return configFile, l, nil
	}
	return "", ServiceLanguageUnknown, fmt.Errorf("service %s not found", serviceName)
}

func ReadServiceConfig(workspaceDir string, serviceName string) (ServiceConfig, error) {
	configFile, language, err := locateServiceConfigPath(workspaceDir, serviceName)
	if err != nil {
		return ServiceConfig{}, err
	}

	rawData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return ServiceConfig{}, fmt.Errorf("failed to read service config: %w", err)
	}

	var data ServiceConfig
	err = yaml.Unmarshal(rawData, &data)
	if err != nil {
		return ServiceConfig{}, fmt.Errorf("failed to read service config: %w", err)
	}
	data.Language = language

	return data, nil
}

func SaveServiceConfig(workspaceDir string, serviceName string, conf ServiceConfig) error {
	svcConfigPath, err := getServiceConfigPathByLang(conf.Language)
	if err != nil {
		return err
	}
	path := filepath.Join(workspaceDir, svcConfigPath, serviceName)
	data, err := yaml.Marshal(&conf)
	if err != nil {
		return fmt.Errorf("failed to create service config: %w", err)
	}

	err = ioutil.WriteFile(filepath.Join(path, ServiceConfigName), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create service config: %w", err)
	}
	return nil
}

func GetServices(workspaceDir string) ([]ServiceInfo, error) {
	svcList := []ServiceInfo{}
	for _, l := range LanguagesList {
		subDir, _ := getServiceConfigPathByLang(l)
		globStr := filepath.Join(workspaceDir, subDir, "*")
		services, err := filepath.Glob(globStr)
		if err != nil {
			return []ServiceInfo{}, err
		}
		for _, f := range services {
			stat, err := os.Stat(f)
			if err != nil {
				return []ServiceInfo{}, err
			}
			if !stat.IsDir() {
				continue
			}
			configFile := filepath.Join(f, ServiceConfigName)
			if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
				continue
			}
			svc := filepath.Base(f)
			svcList = append(svcList, ServiceInfo{
				Name: svc,
				ConfigPath: configFile,
			})
		}
	}
	return svcList, nil
}

func HasService(workspaceDir string, serviceName string) (bool, error) {
	svcList, err := GetServices(workspaceDir)
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

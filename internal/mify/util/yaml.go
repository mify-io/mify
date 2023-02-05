package util

import (
	"os"

	"gopkg.in/yaml.v2"
)

type YAMLData struct {
	path string
}

func NewYAMLData(path string) YAMLData {
	return YAMLData{
		path: path,
	}
}

func (d *YAMLData) ReadFile(out interface{}) error {
	f, err := os.Open(d.path)
	if err != nil {
		return err
	}

	err = yaml.NewDecoder(f).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

func (d *YAMLData) SaveFile(in interface{}) error {
	f, err := os.OpenFile(d.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	err = yaml.NewEncoder(f).Encode(in)
	if err != nil {
		return err
	}
	return nil
}

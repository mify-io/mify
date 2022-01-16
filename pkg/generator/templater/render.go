package templater

import (
	"os"
	"path/filepath"
	"text/template"
)

func RenderTemplate(name string, templateText string, model interface{}, targetPath string) error {
	tmpl := template.New(name)
	var err error
	if tmpl, err = tmpl.Parse(templateText); err != nil {
		return err
	}

	f, err := createWithPath(targetPath)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(f, model); err != nil {
		return err
	}

	return nil
}

// TODO: move to util
func createWithPath(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(path)
}

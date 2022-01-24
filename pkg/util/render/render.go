package render

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func WrapError(text string, err error) error {
	return fmt.Errorf("error while rendering %s: %w", text, err)
}

func RenderTemplate(templateText string, model interface{}, targetPath string) error {
	hash := sha256.Sum256([]byte(templateText))
	name := base64.StdEncoding.EncodeToString(hash[:])

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

func RenderOrSkipTemplate(templateText string, model interface{}, targetPath string) error {
	if _, err := os.Stat(targetPath); err == nil {
		return nil
	}

	return RenderTemplate(templateText, model, targetPath)
}

func createWithPath(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(path)
}

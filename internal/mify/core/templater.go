package core

import (
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type PathTransformerFunc func(context interface{}, path string) (string, error)

const (
	templateExtension = ".tpl"
)

type RenderParams struct {
	// Path to directory with templates tree
	TemplatesPath string

	// Path to save result
	TargetPath string

	// Allows to overwrite the path of file or directory before moving result to target directory
	PathTransformer PathTransformerFunc
}

func renderTemplate(context interface{}, tplPath string, targetPath string) error {
	tmpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0770); err != nil {
		return err
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, context)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(path string, targetPath string) error {
	bytesRead, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(targetPath, bytesRead, 0644)

	if err != nil {
		return err
	}

	return nil
}

func RenderTemplateTree(context interface{}, params RenderParams) error {
	fmt.Printf("Template render: starting... TemplatesPath: %s. TargetPath: %s.\n", params.TemplatesPath, params.TargetPath)

	return filepath.WalkDir(params.TemplatesPath, func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("Template render: visiting %s\n", path)
		if err != nil {
			return err
		}

		destPath := strings.ReplaceAll(path, params.TemplatesPath, "")
		if params.PathTransformer != nil {
			destPath, err = params.PathTransformer(context, destPath)
			if err != nil {
				return err
			}
		}
		destPath = filepath.Join(params.TargetPath, destPath)

		if d.IsDir() {
			fmt.Printf("Template render: found dir %s. Creating: %s\n", path, destPath)
			return os.MkdirAll(destPath, 0755)
		}

		if filepath.Ext(path) == templateExtension {
			filePath := strings.ReplaceAll(destPath, templateExtension, "")
			fmt.Printf("Template render: found tpl %s. Creating: %s\n", path, filePath)
			return renderTemplate(context, path, filePath)
		}

		fmt.Printf("Template render: found file %s. Creating: %s\n", path, destPath)
		copyFile(path, destPath)

		return nil
	})
}

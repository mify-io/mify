package render

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
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

// func(templateText, model, currentText) => migratedText
type MigrationCallback func(string, interface{}, string) (string, error)

type MigrateSettings struct {
	Migrate              bool
	HasUncommitedChanges func() (bool, error)
	Migrations           []MigrationCallback
}

func RenderOrMigrateTemplate(
	templateText string,
	model interface{},
	targetPath string,
	migrateSettings MigrateSettings) error {

	content, err := os.ReadFile(targetPath)
	if errors.Is(err, os.ErrNotExist) {
		return RenderTemplate(templateText, model, targetPath)
	}

	if !migrateSettings.Migrate {
		return nil
	}

	res := string(content)
	for _, migration := range migrateSettings.Migrations {
		res, err = migration(templateText, model, res)
		if err != nil {
			return fmt.Errorf("can't migrate file %s: %w", targetPath, err)
		}
	}

	if res == string(content) {
		return nil
	}

	hasUncommitedChanges, err := migrateSettings.HasUncommitedChanges()
	if err != nil {
		return err
	}

	if hasUncommitedChanges {
		return errors.New("migration can't be applied since file has uncommitted changes")
	}

	err = os.WriteFile(targetPath, []byte(res), os.ModePerm)
	if err != nil {
		return fmt.Errorf("can't write migrated file: %w", err)
	}

	return nil
}

func createWithPath(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(path)
}

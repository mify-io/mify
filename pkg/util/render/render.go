package render

import (
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type renderFlags struct {
	skipExisting bool
	migrate      bool
}

func NewFlags() renderFlags {
	return renderFlags{}
}

func (r renderFlags) SkipExisting() renderFlags {
	r.skipExisting = true
	return r
}

func (r renderFlags) Migrate() renderFlags {
	r.migrate = true
	return r
}

type renderFile struct {
	targetPath   string
	templateName string
	model        any
	flags        renderFlags
}

func NewFile(ctx *gencontext.GenContext, targetPath string) renderFile {
	templateName := path.Base(targetPath) + ".tpl"
	return renderFile{
		targetPath:   targetPath,
		model:        NewDefaultModel(ctx),
		flags:        NewFlags(),
		templateName: templateName,
	}
}

func (f renderFile) SetFlags(flags renderFlags) renderFile {
	f.flags = flags
	return f
}

func (f renderFile) SetModel(model any) renderFile {
	f.model = model
	return f
}

func (f renderFile) SetTemplateName(templateName string) renderFile {
	f.templateName = templateName
	return f
}

func WrapError(text string, err error) error {
	return fmt.Errorf("error while rendering %s: %w", text, err)
}

func RenderMany(templates embed.FS, files ...renderFile) error {
	for _, file := range files {
		data, err := templates.ReadFile(file.templateName)
		if err != nil {
			return WrapError(file.templateName, err)
		}
		if file.flags.skipExisting {
			err = RenderOrSkipTemplate(string(data), file.model, file.targetPath)
		} else {
			err = RenderTemplate(string(data), file.model, file.targetPath)
		}
		if err != nil {
			return WrapError(file.templateName, err)
		}
	}
	return nil
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

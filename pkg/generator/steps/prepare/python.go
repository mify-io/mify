package prepare

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/util/render"
	"go.uber.org/zap"
	"go.uber.org/zap/zapio"
)

const virtualEnvDirName = "venv"

//go:embed requirements.txt.tpl
var requirementsTemplate string

func checkCommand(cmd string) error {
	_, err := exec.LookPath(cmd)
	if err != nil && !errors.Is(err, exec.ErrNotFound) {
		return fmt.Errorf("failed to check if command exists: %w", err)
	}
	if err != nil {
		return fmt.Errorf("command is not found: %s", cmd)
	}
	return nil
}

func prepareVirtualEnv(ctx *gencontext.GenContext, servicesPath string) error {
	if _, err := os.Stat(filepath.Join(servicesPath, virtualEnvDirName, "bin", "activate")); err == nil {
		return nil
	}
	w := &zapio.Writer{Log: ctx.Logger.Desugar(), Level: zap.DebugLevel}
	defer w.Close()

	requirementsPath := filepath.Join(ctx.GetWorkspace().GetPythonServicesAbsPath(), "requirements.txt")
	if err := render.RenderOrSkipTemplate(requirementsTemplate, struct{}{}, requirementsPath); err != nil {
		return render.WrapError("requirements.txt", err)
	}

	pipCmd := exec.Command("python3", "-m", "pip", "install", "--user", "virtualenv")
	pipCmd.Dir = servicesPath
	pipCmd.Stderr = w
	pipCmd.Stdout = w

	err := pipCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to prepare virtual env check if pip installed, for ubuntu install `python3-pip`, error: %w", err)
	}

	venvCmd := exec.Command("python3", "-m", "venv", virtualEnvDirName)
	venvCmd.Dir = servicesPath
	venvCmd.Stderr = w
	venvCmd.Stdout = w

	err = venvCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to prepare virtual env check if venv package installed, for ubuntu install `python3-venv`, error: %w", err)
	}

	pipReqCmd := exec.Command(virtualEnvDirName+"/bin/python3", "-m", "pip", "install", "-r", "requirements.txt")
	pipReqCmd.Dir = servicesPath
	pipReqCmd.Stderr = w
	pipReqCmd.Stdout = w

	err = pipReqCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func preparePython(ctx *gencontext.GenContext) error {
	if ctx.GetMifySchema() == nil {
		return nil
	}
	if ctx.MustGetMifySchema().Language != mifyconfig.ServiceLanguagePython {
		return nil
	}
	if err := checkCommand("python3"); err != nil {
		return err
	}
	servicesDir := ctx.GetWorkspace().GetPythonServicesAbsPath()
	if _, err := os.Stat(servicesDir); errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(servicesDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create py-services dir in workspace: %w", err)
		}
	}
	if err := prepareVirtualEnv(ctx, servicesDir); err != nil {
		return err
	}
	return nil
}

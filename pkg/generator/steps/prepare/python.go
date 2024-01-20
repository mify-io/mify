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

func prepareVirtualEnv(ctx *gencontext.GenContext, requirementsPath string) error {
	servicesPath := ctx.GetWorkspace().GetPythonServicesAbsPath()
	if _, err := os.Stat(filepath.Join(servicesPath, virtualEnvDirName, "bin", "activate")); err == nil {
		return nil
	}
	if !ctx.MustGetMifySchema().Components.Layout.Enabled {
		ctx.Logger.Infof("skipping venv creation without layout enabled")
		return nil
	}
	if _, err := os.Stat(servicesPath); errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(servicesPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create py-services dir in workspace: %w", err)
		}
	}
	w := &zapio.Writer{Log: ctx.Logger.Desugar(), Level: zap.DebugLevel}
	defer w.Close()


	venvCmd := exec.Command("python3", "-m", "venv", virtualEnvDirName)
	venvCmd.Dir = servicesPath
	venvCmd.Stderr = w
	venvCmd.Stdout = w

	err := venvCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to prepare virtual env check if venv package installed, for ubuntu install `python3-venv`, error: %w", err)
	}

	pipReqCmd := exec.Command(virtualEnvDirName+"/bin/python3", "-m", "pip", "install", "-r", requirementsPath)
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
	requirementsPath := filepath.Join(
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetPath().Abs(),
		"requirements.txt",
	)
	if err := render.RenderOrSkipTemplate(requirementsTemplate, struct{}{}, requirementsPath); err != nil {
		return render.WrapError("requirements.txt", err)
	}
	if err := prepareVirtualEnv(ctx, requirementsPath); err != nil {
		return err
	}
	return nil
}

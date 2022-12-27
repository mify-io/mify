package mify

import (
	"os"
	"os/exec"
	"strings"

	"github.com/mify-io/mify/pkg/workspace"
)

func Run(ctx *CliContext, basePath string) error {
	if err := ServiceGenerate(ctx, basePath, workspace.DevRunnerName, false); err != nil {
		return err
	}

	workspace, err := workspace.InitDescription(basePath)
	if err != nil {
		return err
	}

	// TODO: find go in system
	devRunnerMainPath := workspace.GetDevRunnerMainRelPath()
	goServicesPath := workspace.GetGoServicesRelPath()

	devRunnerCmd := exec.Command("go", "run", "."+strings.TrimLeft(devRunnerMainPath, goServicesPath))
	devRunnerCmd.Dir = workspace.GetGoServicesAbsPath()
	devRunnerCmd.Stderr = os.Stderr
	devRunnerCmd.Stdout = os.Stdout

	err = devRunnerCmd.Start()
	if err != nil {
		return err
	}

	err = devRunnerCmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

package util

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

type CommandBuilder struct {
	cmd       *exec.Cmd
	stdoutBuf bytes.Buffer
	stderrBuf bytes.Buffer
}

type CommandResult struct {
	IsOtherError bool
	ExitCode     int
	Stdout       string
	Stderr       string
}

func NewCommand(ctx context.Context, cmdname string, args ...string) *CommandBuilder {
	cmd := exec.CommandContext(ctx, cmdname, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return &CommandBuilder{
		cmd: cmd,
	}
}

func (c *CommandBuilder) StdoutToString() *CommandBuilder {
	c.cmd.Stdout = &c.stdoutBuf
	return c
}

func (c *CommandBuilder) StderrToString() *CommandBuilder {
	c.cmd.Stderr = &c.stderrBuf
	return c
}

func (c *CommandBuilder) DisableStdout() *CommandBuilder {
	c.cmd.Stdout = nil
	return c
}

func (c *CommandBuilder) DisableStderr() *CommandBuilder {
	c.cmd.Stderr = nil
	return c
}

func (c *CommandBuilder) StdinFromString(input string) *CommandBuilder {
	c.cmd.Stdin = bytes.NewReader([]byte(input))
	return c
}

func (c *CommandBuilder) Run() (CommandResult, error) {
	exitcode := 0
	var exitErr error
	err := c.cmd.Run()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			return CommandResult{IsOtherError: true}, err
		}
		exitcode = ee.ExitCode()
		exitErr = fmt.Errorf("command %s returned non-zero exit code: %d", c.cmd.String(), exitcode)
	}
	return CommandResult{
		ExitCode: exitcode,
		Stdout:   c.stdoutBuf.String(),
		Stderr:   c.stderrBuf.String(),
	}, exitErr
}

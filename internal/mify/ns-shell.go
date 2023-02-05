package mify

import (
	"context"
	"encoding/base64"
	"fmt"
	// "os/exec"
	"strings"
	"sync"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mify-io/mify/internal/mify/util"
)

func askForKey(ctx *CliContext) (string, error) {
	sshkey := ""
	prompt := &survey.Input{
		Message: "Paste your public SSH key here",
	}
	err := survey.AskOne(prompt, &sshkey)
	return sshkey, err
}

func decodeKey(ctx *CliContext, key string) (string, bool) {
	sshkeyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", false
	}
	return string(sshkeyBytes), true
}

func getSSHKey(ctx *CliContext) (string, error) {
	for {
		var sshkey string
		if len(ctx.Config.SSHPublicKey) > 0 {
			if key, ok := decodeKey(ctx, ctx.Config.SSHPublicKey); ok {
				sshkey = key
			}
		}
		var err error
		if len(sshkey) == 0 {
			sshkey, err = askForKey(ctx)
			if err != nil {
				return "", err
			}
		}
		res, err := util.NewCommand(ctx.GetCtx(), "ssh-keygen", "-l", "-f", "-").
			DisableStdout().
			StdinFromString(sshkey).Run()
		if err != nil && res.IsOtherError {
			return "", fmt.Errorf("failed to verify ssh key: %w", err)
		}
		if res.ExitCode == 255 {
			ctx.Config.SSHPublicKey = ""
			continue
		}

		ctx.Config.SSHPublicKey = base64.StdEncoding.EncodeToString([]byte(sshkey))
		return sshkey, nil
	}
}

func NsShell(ctx *CliContext, env string, forwardProxy string, listenPort string) error {
	sshkey, err := getSSHKey(ctx)
	if err != nil {
		return err
	}
	k8sContext, err := getKubeContextName(ctx, env)
	if err != nil {
		return err
	}
	goCtx, cancel := context.WithCancel(ctx.GetCtx())
	defer cancel()

	res, err := util.NewCommand(goCtx,
		"kubectl", "--context="+k8sContext,
		"get", "pods", "-l", "app=ns-shell",
		"-o", "custom-columns=:metadata.name", "--no-headers").StdoutToString().Run()
	if err != nil {
		return err
	}
	podName := strings.TrimSpace(res.Stdout)

	res, err = util.NewCommand(goCtx,
		"kubectl", "--context="+k8sContext,
		"exec", podName, "--",
		"/bin/sh", "-c",
		fmt.Sprintf(
			"grep -q -F \"%s\" ~/.ssh/authorized_keys 2>/dev/null || echo \"%s\" >> ~/.ssh/authorized_keys",
			sshkey, sshkey),
	).StdoutToString().Run()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var fwdRes util.CommandResult
	var fwdErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		fwdRes, fwdErr = util.NewCommand(goCtx,
			"kubectl", "--context="+k8sContext,
			"port-forward", podName, listenPort+":22",
		).DisableStderr().DisableStdout().Run()
	}()

	attemps := 10
	for {
		res, err := util.NewCommand(goCtx, "ssh", "-q", "-p", listenPort, "user@localhost", "exit").Run()
		if err == nil {
			break
		}
		if res.ExitCode == 255 {
			if _, err := util.NewCommand(goCtx, "ssh-keygen", "-R", "[localhost]:"+listenPort).DisableStderr().Run(); err != nil {
				return err
			}
		}
		if err != nil && attemps == 0 {
			return fmt.Errorf("failed to connect to pod via ssh, please check your public key configuration: %w", err)
		}
		attemps -= 1
		time.Sleep(100 * time.Millisecond)
	}

	sshArgs := []string{"-p", listenPort}
	if forwardProxy != "" {
		sshArgs = append(sshArgs, "-L", forwardProxy)
	}

	fmt.Printf("You can now connect to pod via ssh, run this command to create additional sessions:\n")
	fmt.Printf("$ ssh -p %s user@localhost\n", listenPort)
	fmt.Printf("============================\n")

	sshArgs = append(sshArgs, "user@localhost")
	res, err = util.NewCommand(goCtx, "ssh", sshArgs...).Run()
	if err != nil && res.IsOtherError {
		return err
	}

	cancel()
	wg.Wait()
	if fwdErr != nil && fwdRes.IsOtherError {
		return fwdErr
	}
	return nil
}

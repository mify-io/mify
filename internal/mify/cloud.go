package mify

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mify-io/mify/internal/mify/util"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace"
)

const CLOUD_URL = "https://cloud.mify.io"

func CloudInit(ctx *CliContext, projectName string, env string) error {
	_, err := workspace.InitDescription(ctx.WorkspacePath)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	ctx.Logger.Printf("Please visit %s to receive token and paste it here:", CLOUD_URL)

	token, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read token from stdin: %w", err)
	}
	token = strings.TrimSpace(token)
	accessToken, err := getAccessToken(ctx, token)
	if err != nil {
		return fmt.Errorf("token validation error: %w", err)
	}
	ctx.Config.APIToken = token
	err = SaveConfig(ctx.Config)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	if len(projectName) == 0 {
		projectName = ctx.workspaceDescription.Name
		ctx.Logger.Printf("Your project name (default: %s):", projectName)
		newName, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read project name from stdin: %w", err)
		}
		newName = strings.TrimSpace(newName)
		if len(newName) > 0 {
			projectName = newName
		}
	}
	if len(env) == 0 {
		env = "stage"
		ctx.Logger.Printf(`Project environment ("stage" or "prod", default: "stage"):`)
		newEnv, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read environment name from stdin: %w", err)
		}
		newEnv = strings.TrimSpace(newEnv)
		if len(newEnv) > 0 {
			env = newEnv
		}
	}
	ctx.workspaceDescription.Config.ProjectName = projectName
	ctx.workspaceDescription.Config.Environments = util.StringSetAppend(ctx.workspaceDescription.Config.Environments, env)
	err = mifyconfig.SaveWorkspaceConfig(ctx.WorkspacePath, ctx.workspaceDescription.Config)
	if err != nil {
		return fmt.Errorf("failed to update workspace config: %w", err)
	}

	err = registerProject(ctx, projectName, env, accessToken)
	if err != nil {
		return fmt.Errorf("failed to register project: %w", err)
	}
	fmt.Println("Successfully registered project! Now you can deploy services via `mify cloud deploy`.")

	return nil
}

func getAccessToken(ctx *CliContext, token string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/auth/token/service", CLOUD_URL)
	var reqData struct {
		RefreshToken string `json:"refresh_token"`
	}
	var respData struct {
		AccessToken string `json:"access_token"`
	}
	reqData.RefreshToken = token
	client := resty.New()
	resp, err := client.R().SetBody(reqData).SetResult(&respData).Post(endpoint)
	if err != nil {
		return "", fmt.Errorf("request to get token failed: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("request to get token error: %s", resp.Status())
	}
	return respData.AccessToken, nil

}

func registerProject(ctx *CliContext, projectName string, environment string, accessToken string) error {
	endpoint := fmt.Sprintf("%s/api/projects/register", CLOUD_URL)
	var reqData struct {
		Name        string `json:"name"`
		Environment string `json:"environment"`
	}
	reqData.Name = projectName
	reqData.Environment = environment
	client := resty.New()
	resp, err := client.R().SetAuthToken(accessToken).SetBody(reqData).Post(endpoint)
	if err != nil {
		return fmt.Errorf("request to get token failed: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request to get token error: %s", resp.Status())
	}
	return nil
}

func CloudUpdateKubeconfig(ctx *CliContext) error {
	return nil
}

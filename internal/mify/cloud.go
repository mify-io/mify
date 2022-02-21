package mify

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mify-io/mify/internal/mify/util"
	"github.com/mify-io/mify/pkg/cloudconfig"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

const CLOUD_URL = "https://cloud.mify.io"

func CloudInit(ctx *CliContext, projectName string, env string) error {
	if ctx.Config.APIToken == "" {
		if err := obtainApiToken(ctx); err != nil {
			return err
		}
	}

	accessToken, err := resolveAccessToken(ctx)
	if err != nil {
		return err
	}

	if len(projectName) == 0 {
		projectName = ctx.workspaceDescription.Name
		newName, err := ctx.UserInput.AskInput("Your project name (default: %s):", projectName)
		if err != nil {
			return fmt.Errorf("failed to read project name from stdin: %w", err)
		}
		if len(newName) > 0 {
			projectName = newName
		}
	}
	if len(env) == 0 {
		env = "stage"
		newEnv, err := ctx.UserInput.AskInput(`Project environment ("stage" or "prod", default: "stage"):`)
		if err != nil {
			return fmt.Errorf("failed to read environment name from stdin: %w", err)
		}
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

	err = initCloudConfigs(ctx)
	if err != nil {
		return fmt.Errorf("failed to init cloud configs: %w", err)
	}

	return nil
}

func obtainApiToken(ctx *CliContext) error {
	token, err := ctx.UserInput.AskInput("Please visit %s to receive token and paste it here:", CLOUD_URL)

	if err != nil {
		return fmt.Errorf("failed to read token from stdin: %w", err)
	}
	ctx.Config.APIToken = strings.TrimSpace(token)

	err = SaveConfig(ctx.Config)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

func resolveAccessToken(ctx *CliContext) (string, error) {
	accessToken, err := getAccessToken(ctx, ctx.Config.APIToken)
	if err != nil {
		return "", fmt.Errorf("token validation error: %w", err)
	}

	return accessToken, nil
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

func initCloudConfigs(ctx *CliContext) error {
	for _, service := range ctx.workspaceDescription.GetApiServices() {
		path := ctx.MustGetWorkspaceDescription().GetCloudSchemaAbsPath(service)
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			continue
		}

		config := &cloudconfig.ServiceCloudConfig{}
		err := config.WriteToFile(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func CloudUpdateKubeconfig(ctx *CliContext) error {
	return nil
}

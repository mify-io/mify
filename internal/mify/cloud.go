package mify

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"

	"github.com/mify-io/mify/pkg/cloudconfig"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace/mutators/cloud"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func CloudInit(ctx *CliContext) error {
	if ctx.Config.APIToken == "" {
		if err := obtainApiToken(ctx); err != nil {
			return err
		}
	}

	accessToken, err := resolveAccessToken(ctx)
	if err != nil {
		return err
	}

	orgsList, err := getOrganizations(ctx, accessToken)
	if err != nil {
		return err
	}

	if len(orgsList.Organizations) == 0 {
		return fmt.Errorf(
			"You have no organizations registered in Mify Cloud, please visit %s to register.",
			cloudconfig.GetCloudUrl(),
		)
	}
	org := orgsList.Organizations[0]
	orgID := org.Id
	if len(orgsList.Organizations) > 1 {
		names := lo.Map(orgsList.Organizations, func(org organization, _ int) string {
			return org.Name
		})
		var selectedName string
		err := survey.AskOne(&survey.Select{
			Message: "Choose organization to link to this workspace:",
			Options: names,
			Default: names[0],
		}, &selectedName)
		if err != nil {
			return err
		}
		orgs := lo.Filter(orgsList.Organizations, func(org organization, _ int) bool {
			return org.Name == selectedName
		})
		org = orgs[0]
		orgID = org.Id
	}

	projectName := ctx.workspaceDescription.Config.ProjectName
	curProjects := lo.Filter(org.Projects, func(p project, _ int) bool {
		return p.Name == projectName
	})

	if len(curProjects) == 0 {
		curProject := org.Projects[0]
		projectName = curProject.Name
		curProjects = lo.Filter(org.Projects, func(p project, _ int) bool {
			return p.Name == projectName
		})
		projectNames := lo.Uniq(lo.Map(org.Projects, func(p project, _ int) string {
			return p.Name
		}))
		if len(projectNames) > 1 {
			err := survey.AskOne(&survey.Select{
				Message: "Choose project to link to this workspace:",
				Options: projectNames,
				Default: projectNames[0],
			}, &projectName)
			if err != nil {
				return err
			}
			curProjects = lo.Filter(org.Projects, func(p project, _ int) bool {
				return p.Name == projectName
			})
			curProject = curProjects[0]
		}
	}
	envs := lo.FilterMap(curProjects, func(p project, _ int) (string, bool) {
		if p.Name == projectName {
			return p.Environment, true
		}
		return "", false
	})
	ctx.workspaceDescription.Config.OrganizationID = orgID
	ctx.workspaceDescription.Config.ProjectName = projectName
	ctx.workspaceDescription.Config.Environments = envs
	err = mifyconfig.SaveWorkspaceConfig(ctx.WorkspacePath, ctx.workspaceDescription.Config)
	if err != nil {
		return fmt.Errorf("failed to update workspace config: %w", err)
	}

	fmt.Println("Successfully registered project! Now you can deploy services via `mify cloud deploy`.")

	err = initCloudConfigs(ctx)
	if err != nil {
		return fmt.Errorf("failed to init cloud configs: %w", err)
	}

	if err := cloud.Init(ctx.mutatorContext); err != nil {
		return err
	}
	for _, env := range envs {
		if err := CloudUpdateKubeconfig(ctx, env); err != nil {
			return err
		}
	}

	return nil
}

func obtainApiToken(ctx *CliContext) error {
	token, err := ctx.UserInput.AskInput(
		"Please visit %s to receive token and paste it here:",
		cloudconfig.GetCloudUrl())

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
	endpoint := fmt.Sprintf("%s/auth/token/service", cloudconfig.GetCloudApiURL())
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

type project struct {
	Environment string `json:"environment"`
	Id          string `json:"id"`
	Name        string `json:"name"`
}
type organization struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Projects []project `json:"projects"`
}
type orgsResponse struct {
	Organizations []organization `json:"organizations"`
}

func getOrganizations(ctx *CliContext, accessToken string) (orgsResponse, error) {
	var out orgsResponse
	endpoint := fmt.Sprintf("%s/orgs", cloudconfig.GetCloudApiURL())
	client := resty.New()
	resp, err := client.R().SetAuthToken(accessToken).SetResult(&out).Get(endpoint)
	if err != nil {
		return orgsResponse{}, fmt.Errorf("request to get token failed: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return orgsResponse{}, fmt.Errorf("request to get token error: %s", resp.Status())
	}
	return out, nil
}

func initCloudConfigs(ctx *CliContext) error {
	for _, service := range ctx.workspaceDescription.GetApiServices() {
		path := ctx.MustGetWorkspaceDescription().GetCloudSchemaAbsPath(service)
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			continue
		}

		_, err := os.Create(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func findKubeConfig() (string, error) {
	env := os.Getenv("KUBECONFIG")
	if env != "" {
		return env, nil
	}
	path, err := homedir.Expand("~/.kube/config")
	if err != nil {
		return "", err
	}
	return path, nil
}

type kubeconfigResponse struct {
	ServerAddress    string `json:"server_address"`
	ServerCertficate string `json:"server_certficate"`
	ServiceAccount   string `json:"service_account"`
	Token            string `json:"token"`
}

func getKubeconfigData(
	ctx *CliContext, projectName string,
	environment string, accessToken string) (kubeconfigResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/kubeconfig", cloudconfig.GetCloudApiURL())
	var reqData struct {
		Name        string `json:"name"`
		Environment string `json:"environment"`
	}
	reqData.Name = projectName
	reqData.Environment = environment
	client := resty.New()
	var result kubeconfigResponse
	resp, err := client.R().SetAuthToken(accessToken).SetBody(reqData).SetResult(&result).Post(endpoint)
	if err != nil {
		return kubeconfigResponse{}, fmt.Errorf("request to get token failed: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return kubeconfigResponse{}, fmt.Errorf("request to get token error: %s", resp.Status())
	}
	return result, nil

}

func CloudUpdateKubeconfig(ctx *CliContext, environment string) error {
	accessToken, err := resolveAccessToken(ctx)
	if err != nil {
		return err
	}
	wspc := ctx.MustGetWorkspaceDescription()
	data, err := getKubeconfigData(ctx, wspc.Name, environment, accessToken)
	if err != nil {
		return fmt.Errorf("failed to register project: %w", err)
	}
	cert, err := base64.StdEncoding.DecodeString(data.ServerCertficate)
	if err != nil {
		return err
	}
	token, err := base64.StdEncoding.DecodeString(data.Token)
	if err != nil {
		return err
	}

	kubeConfigPath, err := findKubeConfig()
	if err != nil {
		return err
	}

	var kubeConfig *api.Config
	if _, err := os.Stat(kubeConfigPath); errors.Is(err, os.ErrNotExist) {
		kubeConfig = api.NewConfig()
	} else {
		kubeConfig, err = clientcmd.LoadFromFile(kubeConfigPath)
		if err != nil {
			return err
		}
	}

	clusterName := "mifykube-" + environment
	cluster := api.NewCluster()
	cluster.CertificateAuthorityData = cert
	cluster.Server = data.ServerAddress

	context := api.NewContext()
	context.Cluster = clusterName
	context.Namespace = wspc.Name + "-" + environment
	context.AuthInfo = data.ServiceAccount

	user := api.NewAuthInfo()
	user.Token = string(token)

	contextName := data.ServiceAccount + "@" + clusterName
	kubeConfig.Clusters[clusterName] = cluster
	kubeConfig.Contexts[contextName] = context
	kubeConfig.AuthInfos[data.ServiceAccount] = user
	kubeConfig.CurrentContext = contextName

	err = clientcmd.WriteToFile(*kubeConfig, kubeConfigPath)
	if err != nil {
		return err
	}
	return nil
}

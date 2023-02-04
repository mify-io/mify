package auth

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/mify-io/mify/pkg/cloudconfig"
)

func ResolveAccessToken(apiToken string) (string, error) {
	accessToken, err := getAccessToken(apiToken)
	if err != nil {
		return "", fmt.Errorf("token validation error: %w", err)
	}

	return accessToken, nil
}

func getAccessToken(token string) (string, error) {
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

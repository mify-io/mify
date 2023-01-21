package stats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type SendStatsReq struct {
	Events []Event `json:"events"`
}

func SendStats(mifyStatsApiUrl string, apiToken string, events []Event) error {
	data, err := json.Marshal(SendStatsReq{Events: events})
	if err != nil {
		return err
	}

	client := resty.New()
	client.SetTimeout(1 * time.Second)
	endpoint := fmt.Sprintf("%s/events/cli", mifyStatsApiUrl)
	resp, err := client.R().SetAuthToken(apiToken).SetBody(data).Put(endpoint)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("status code: %s, body: %s", resp.Status(), resp.Body())
	}

	return nil
}

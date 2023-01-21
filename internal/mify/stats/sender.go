package stats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
)

type SendStatsReq struct {
	Events []Event `json:"events"`
}

const SendStatsThreashhold = 10

func sendStatsToServer(url string, apiToken string, data []byte) error {
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	endpoint := fmt.Sprintf("%s/events/cli", url)
	resp, err := client.R().SetAuthToken(apiToken).SetBody(data).Put(endpoint)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("status code: %s, body: %s", resp.Status(), resp.Body())
	}

	return nil
}

// Send all stats saved in queue file (only if there are more then SendStatsThreashhold events)
func (s *Collector) MaybeSendStats() error {
	s.statsQueueFileMutex.Lock()
	defer s.statsQueueFileMutex.Unlock()

	f, err := os.ReadFile(s.statsQueueFile)
	if err != nil {
		return fmt.Errorf("can't open stat event queue file: %w", err)
	}

	lines := lo.WithoutEmpty(strings.Split(string(f), "\n"))
	if len(lines) < SendStatsThreashhold {
		return nil
	}

	s.logger.Printf("about to send statistics to server: %v events", len(lines))

	events := lo.FilterMap(lines, func(line string, _ int) (Event, bool) {
		var ev Event
		err := json.Unmarshal([]byte(line), &ev)
		if err != nil {
			s.logger.Printf("stats queue file is corrupted, skipping saved event: %s", err)
			return ev, false
		}
		return ev, true
	})

	data, err := json.Marshal(SendStatsReq{Events: events})
	if err != nil {
		return err
	}

	err = sendStatsToServer(s.apiUrl, s.apiToken, data)
	if err != nil {
		return err
	}

	err = os.Remove(s.statsQueueFile)
	if err != nil {
		return err
	}

	return nil
}

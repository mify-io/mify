package stats

import (
	"context"
	"log"
	"runtime"

	"github.com/mify-io/mify/pkg/cloudconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Collector struct {
	logger     *log.Logger
	isEnabled  bool
	apiUrl     string
	instanceID string
	apiToken   string
	ctx        context.Context
}

type CmdFlag struct {
	Name         string `json:"name"`
	Value        string `json:"value,omitempty"`
	DefaultValue string `json:"default_value"`
	IsChanged    bool   `json:"is_changed"`
}

type EventPayload struct {
	Flags []CmdFlag `json:"flags"`
	Args  []string  `json:"args"`
}

type Event struct {
	Name       string       `json:"name"`
	OS         string       `json:"os"`
	Arch       string       `json:"arch"`
	InstanceID string       `json:"mify_instance_id"`
	Payload    EventPayload `json:"payload"`
}

func NewCollector(
	ctx context.Context,
	logger *log.Logger,
	isEnabled bool, instanceID string, apiToken string) *Collector {
	return &Collector{
		ctx:        ctx,
		logger:     logger,
		isEnabled:  isEnabled,
		instanceID: instanceID,
		apiToken:   apiToken,
		apiUrl:     cloudconfig.GetCloudStatsURL(),
	}
}

func (s *Collector) LogEvent(name string, cmd *cobra.Command) {
	if !s.isEnabled {
		return
	}
	os := runtime.GOOS
	arch := runtime.GOARCH
	flags := []CmdFlag{}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		flags = append(flags, CmdFlag{
			Name:         f.Name,
			Value:        f.Value.String(),
			DefaultValue: f.DefValue,
			IsChanged:    f.Changed,
		})
	})
	event := Event{
		Name:       name,
		OS:         os,
		Arch:       arch,
		InstanceID: s.instanceID,
		Payload: EventPayload{
			Flags: flags,
			Args:  cmd.Flags().Args(),
		},
	}
	s.logger.Printf("log\n")
	s.logger.Printf("event: %+v\n", event)
	// go func() {
	// resty
	// }
	s.logger.Printf("done\n\n")
}

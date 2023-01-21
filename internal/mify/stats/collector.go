package stats

import (
	"context"
	"encoding/json"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mify-io/mify/pkg/cloudconfig"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Collector struct {
	logger        *log.Logger
	isEnabled     bool
	apiUrl        string
	instanceID    string
	workspaceName string
	projectName   string
	mifyVersion   string
	apiToken      string
	ctx           context.Context
}

type CmdCommand struct {
	Name  string      `json:"name"`
	Flags []CmdFlag   `json:"flags"`
	Args  []string    `json:"args"`
	Child *CmdCommand `json:"child"`
}

type CmdFlag struct {
	Name         string `json:"name"`
	Value        string `json:"value,omitempty"`
	DefaultValue string `json:"default_value"`
	IsChanged    bool   `json:"is_changed"`
}

type RunPayload struct {
	CmdInfo []CmdCommand `json:"cmd_info"`
}

type Event struct {
	Id             string `json:"id"`
	UserTime       string `json:"user_time"`
	Name           string `json:"name"`
	MifyVersion    string `json:"mify_version"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	MifyInstanceID string `json:"mify_instance_id"`
	WorkspaceName  string `json:"workspace_name"`
	ProjectName    string `json:"project_name"`
	Payload        string `json:"payload"`
}

func NewCollector(
	ctx context.Context,
	logger *log.Logger,
	isEnabled bool,
	instanceID string,
	workspaceName string,
	projectName string,
	mifyVersion string,
	apiToken string) *Collector {
	return &Collector{
		ctx:           ctx,
		logger:        logger,
		isEnabled:     isEnabled,
		instanceID:    instanceID,
		workspaceName: workspaceName,
		projectName:   projectName,
		mifyVersion:   mifyVersion,
		apiToken:      apiToken,
		apiUrl:        cloudconfig.GetStatsApiUrl(),
	}
}

func WrapCmdCommand(command *cobra.Command) *CmdCommand {
	var flags []CmdFlag
	command.Flags().VisitAll(func(f *pflag.Flag) {
		flags = append(flags, CmdFlag{
			Name:         f.Name,
			Value:        f.Value.String(),
			DefaultValue: f.DefValue,
			IsChanged:    f.Changed,
		})
	})

	return &CmdCommand{
		Name:  command.Name(),
		Flags: flags,
		Args:  command.ValidArgs,
	}
}

func Trim(s string, size int) string {
	if len(s) <= size {
		return s
	}

	return s[0:size]
}

func (s *Collector) LogCobraCommandExecuted(cmd *cobra.Command) {
	if !s.isEnabled || s.mifyVersion == "" {
		return
	}

	if strings.HasPrefix(cmd.Name(), "__") {
		// ignore __complete and other
		return
	}

	os := runtime.GOOS
	arch := runtime.GOARCH

	var commandInfo []CmdCommand
	commandInfo = append(commandInfo, *WrapCmdCommand(cmd))
	cmd.VisitParents(func(c *cobra.Command) {
		commandInfo = append(commandInfo, *WrapCmdCommand(c))
	})

	data, err := json.Marshal(RunPayload{
		CmdInfo: lo.Reverse(commandInfo),
	})
	if err != nil {
		panic(err)
	}

	event := Event{
		Id:             uuid.New().String(),
		UserTime:       time.Now().UTC().Format(time.RFC3339),
		Name:           "run",
		MifyVersion:    s.mifyVersion,
		OS:             Trim(os, 128),
		Arch:           Trim(arch, 128),
		MifyInstanceID: s.instanceID,
		WorkspaceName:  Trim(s.workspaceName, 64),
		ProjectName:    Trim(s.projectName, 64),
		Payload:        string(data),
	}

	err = SendStats(s.apiUrl, s.apiToken, []Event{event})
	if err != nil {
		s.logger.Printf("Warn: can't send usage statistics to mify.io: %s", err)
	}
}

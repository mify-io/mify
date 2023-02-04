package mify

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"github.com/mify-io/mify/internal/mify/stats"
	"github.com/mify-io/mify/internal/mify/userinput"
	"github.com/mify-io/mify/pkg/workspace"
	"github.com/mify-io/mify/pkg/workspace/mutators"
	"github.com/spf13/viper"
)

type Config struct {
	APIToken          string `mapstructure:"MIFY_API_TOKEN"`
	DisableUsageStats bool   `mapstructure:"MIFY_DISABLE_USAGE_STATS"`
	InstanceID        string `mapstructure:"MIFY_INSTANCE_ID"`
	SSHPublicKey      string `mapstructure:"SSH_PUBLIC_KEY"`
}

func (c Config) Equal(newConfig Config) bool {
	if c.APIToken != newConfig.APIToken {
		return false
	}
	if c.DisableUsageStats != newConfig.DisableUsageStats {
		return false
	}
	if c.InstanceID != newConfig.InstanceID {
		return false
	}
	if c.SSHPublicKey != newConfig.SSHPublicKey {
		return false
	}
	return true
}

func GetConfigDirectory() string {
	return filepath.Join(xdg.ConfigHome, "mify")
}

func NewDefaultConfig() Config {
	instanceID := uuid.New().String()
	viper.SetDefault("MIFY_API_TOKEN", "")
	viper.SetDefault("MIFY_DISABLE_USAGE_STATS", false)
	viper.SetDefault("MIFY_INSTANCE_ID", instanceID)
	viper.SetDefault("SSH_PUBLIC_KEY", "")
	return Config{
		APIToken:          "",
		DisableUsageStats: false,
		InstanceID:        instanceID,
		SSHPublicKey: "",
	}
}

func SaveConfig(config Config) error {
	viper.Set("MIFY_API_TOKEN", config.APIToken)
	viper.Set("MIFY_DISABLE_USAGE_STATS", config.DisableUsageStats)
	viper.Set("MIFY_INSTANCE_ID", config.InstanceID)
	viper.Set("SSH_PUBLIC_KEY", config.SSHPublicKey)

	err := os.MkdirAll(GetConfigDirectory(), 0755)
	if err != nil {
		return err
	}
	configPath := filepath.Join(GetConfigDirectory(), "config.yaml")
	return viper.WriteConfigAs(configPath)
}

func UpdateConfig(config Config) error {
	configPath := filepath.Join(GetConfigDirectory(), "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return SaveConfig(config)
	}
	oldConfig := NewDefaultConfig()
	if err := viper.Unmarshal(&oldConfig); err != nil {
		return fmt.Errorf("failed to read config: %s", err)
	}
	if !oldConfig.Equal(config) {
		return SaveConfig(config)
	}
	return nil
}

type CliContext struct {
	Logger         *log.Logger
	Ctx            context.Context
	Cancel         context.CancelFunc
	Config         Config
	WorkspacePath  string
	MifyVersion    string
	IsVerbose      bool
	UserInput      userinput.UserInput
	StatsCollector *stats.Collector

	workspaceDescription *workspace.Description
	mutatorContext       *mutators.MutatorContext
}

func NewContext(config Config, workspacePath string, isVerbose bool, mifyVersion string) *CliContext {
	ctx, cancel := context.WithCancel(context.Background())
	logger := log.New(os.Stdout, "", 0)
	return &CliContext{
		Logger:        logger,
		Ctx:           ctx,
		Cancel:        cancel,
		Config:        config,
		WorkspacePath: workspacePath,
		MifyVersion:   mifyVersion,
		IsVerbose:     isVerbose,
	}
}

func (c *CliContext) InitStatsCollector(statsQueueFile string) {
	workspaceName := ""
	projectName := ""
	if c.workspaceDescription != nil {
		workspaceName = c.workspaceDescription.Config.WorkspaceName
		projectName = c.workspaceDescription.Config.ProjectName
	}

	c.StatsCollector = stats.NewCollector(
		c.Ctx,
		c.Logger,
		!c.Config.DisableUsageStats || c.MifyVersion == "", // mifyVersion == "" => non release build
		c.Config.InstanceID,
		workspaceName,
		projectName,
		c.MifyVersion,
		c.Config.APIToken,
		statsQueueFile)
}

func (c CliContext) GetCtx() context.Context {
	return c.Ctx
}

func (c *CliContext) LoadWorkspace() error {
	res, err := workspace.InitDescription(c.WorkspacePath)
	if err != nil {
		return err
	}

	c.workspaceDescription = &res
	c.WorkspacePath = c.workspaceDescription.BasePath

	c.mutatorContext = mutators.NewMutatorContext(c.Ctx, c.Logger, c.workspaceDescription)

	return nil
}

func (c *CliContext) MustGetWorkspaceDescription() *workspace.Description {
	if c.workspaceDescription == nil {
		panic("missed workspaceDescription")
	}

	return c.workspaceDescription
}

func (c *CliContext) GetWorkspaceDescription() *workspace.Description {
	return c.workspaceDescription
}

func (c *CliContext) MustGetMutatorContext() *mutators.MutatorContext {
	if c.mutatorContext == nil {
		panic("missed mutatorContext")
	}

	return c.mutatorContext
}

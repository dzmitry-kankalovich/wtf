package awscosts

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable = true
	defaultTitle     = "AWS Costs"
	defaultTopN      = -1
)

type creds struct {
	profile string `help:"AWS profile" default:"default" optional:"true"`
}

type Settings struct {
	*cfg.Common

	account string `help:"AWS account ID" default:"" optional:"false"`
	alias   string `help:"A human-readable alias for AWS Account" default:"<AccountID>" optional:"true"`
	topN    int    `help:"Display top N services by cost. Specify N, otherwise no per-service breakdown done" default:"-1" optional:"true"`
	creds
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := &Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),
	}

	settings.account = ymlConfig.UString("account", "")
	settings.alias = ymlConfig.UString("alias", settings.account)
	settings.topN = ymlConfig.UInt("topN", defaultTopN)

	settings.creds.profile = ymlConfig.UString("creds.profile", "default")

	return settings
}

package githubtrending

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable     = true
	defaultTitle         = "GitHub Trending"
	defaultPeriod        = "day"
	defaultNumberOfRepos = 10
	defaultMinStars      = 0
	defaultUseScraper    = false
)

type colors struct {
	repo  string `help:"Color to use for repository name." optional:"true" default:"green"`
	lang  string `help:"Color to use for repository language." optional:"true" default:"orange"`
	stars string `help:"Color to use for stars count." optional:"true" default:"yellow"`
}

type Settings struct {
	*cfg.Common

	colors

	limit       int      `help:"Number of repositories to be displayed" default:"10" optional:"true"`
	langs       []string `help:"Filter results by programming languages" optional:"true"`
	minStars    int      `help:"Minimum amount of stars" default:"0" optional:"true"`
	spokenLangs []string `help:"Filter results by spoken languages. Works only with scraper enabled, otherwise ignored" optional:"true"`
	period      string   `help:"Time span for the trending repositories. Options are 'day', 'week' and 'month'" default:"day" optional:"true"`
	useScraper  bool     `help:"Use more advanced, but error-prone scraping of https://github.com/trending instead of GitHub API. Allows filtering by spoken language" default:"false" optional:"true"`
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := &Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		limit:       ymlConfig.UInt("limit", defaultNumberOfRepos),
		langs:       cfg.ParseAsMapOrList(ymlConfig, "langs"),
		minStars:    ymlConfig.UInt("minStars", defaultMinStars),
		spokenLangs: cfg.ParseAsMapOrList(ymlConfig, "spokenLangs"),
		useScraper:  ymlConfig.UBool("useScraper", defaultUseScraper),
		period:      ymlConfig.UString("period", defaultPeriod),
	}

	settings.colors.repo = ymlConfig.UString("colors.repo", "green")
	settings.colors.lang = ymlConfig.UString("colors.lang", "orange")
	settings.colors.stars = ymlConfig.UString("colors.lang", "yellow")

	return settings
}

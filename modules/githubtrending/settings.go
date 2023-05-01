package githubtrending

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable     = true
	defaultTitle         = "GitHub Trending"
	defaultSince         = "daily"
	defaultNumberOfRepos = 10
	defaultMinStars      = 0
	defaultUseScraper    = true
)

type colors struct {
	repo  string `help:"Color to use for repository name." optional:"true" default:"green"`
	lang  string `help:"Color to use for repository language." optional:"true" default:"orange"`
	stars string `help:"Color to use for stars count." optional:"true" default:"yellow"`
}

type icons struct {
	stars     string `help:"An icon char to use for stars" optional:"true" default:"⭐"`
	starsDiff string `help:"An icon char to use for stars gained today" optional:"true" default:"⬆️"`
}

type Settings struct {
	*cfg.Common

	colors
	icons

	limit       int      `help:"Number of repositories to be displayed" default:"10" optional:"true"`
	langs       []string `help:"Filter results by programming languages" optional:"true"`
	minStars    int      `help:"Minimum amount of stars" default:"0" optional:"true"`
	spokenLangs []string `help:"Filter results by spoken languages. Works only with scraper enabled, otherwise ignored" optional:"true"`
	since       string   `help:"Time period for the trending repositories. Options are 'daily', 'weekly' and 'monthly'" default:"daily" optional:"true"`
	useScraper  bool     `help:"Use more advanced, but error-prone scraping of https://github.com/trending instead of GitHub API. Allows filtering by spoken language" default:"true" optional:"true"`
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := &Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		limit:       ymlConfig.UInt("limit", defaultNumberOfRepos),
		langs:       cfg.ParseAsMapOrList(ymlConfig, "langs"),
		minStars:    ymlConfig.UInt("minStars", defaultMinStars),
		spokenLangs: cfg.ParseAsMapOrList(ymlConfig, "spokenLangs"),
		useScraper:  ymlConfig.UBool("useScraper", defaultUseScraper),
		since:       ymlConfig.UString("since", defaultSince),
	}

	settings.colors.repo = ymlConfig.UString("colors.repo", "green")
	settings.colors.lang = ymlConfig.UString("colors.lang", "orange")
	settings.colors.stars = ymlConfig.UString("colors.lang", "yellow")

	settings.icons.stars = ymlConfig.UString("icons.stars", "⭐")

	// Known problem with this character in iTerm2: https://gitlab.com/gnachman/iterm2/-/issues/8735
	// Fix for iTerm2: add a space immediately after the character
	// Alternatively set your own character via settings.icons.starsDiff
	settings.icons.starsDiff = ymlConfig.UString("icons.starsDiff", "⬆️ ")

	return settings
}

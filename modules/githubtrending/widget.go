package githubtrending

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
	"math"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	ghb "github.com/google/go-github/v32/github"
)

type ShowType int

const (
	ShowName ShowType = iota
	ShowNameDesc
	ShowNameDescLang
	ShowNameDescLangStars
)

type Widget struct {
	view.ScrollableWidget

	repos     []Repo
	settings  *Settings
	err       error
	client    *ghb.Client
	collector *colly.Collector
	showType  ShowType
}

func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := &Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.Common),
		settings:         settings,
		showType:         ShowNameDescLangStars,
	}

	if settings.useScraper {
		widget.collector = colly.NewCollector(
			colly.AllowURLRevisit(),
		)
	} else {
		widget.client = ghb.NewClient(nil)
	}

	widget.SetRenderFunction(widget.Render)
	widget.initializeKeyboardControls()

	return widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	var repos []Repo
	var err error

	if widget.settings.useScraper {
		repos, err = ScrapeRepos(widget.collector, widget.settings)
	} else {
		repos, err = FetchRepos(widget.client, widget.settings)
	}

	if err != nil {
		widget.err = err
		widget.repos = nil
		widget.SetItemCount(0)
	} else {
		widget.repos = repos
		widget.SetItemCount(len(repos))
	}

	widget.Render()
}

func (widget *Widget) Render() {
	widget.Redraw(widget.content)
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() (string, string, bool) {
	title := fmt.Sprintf("%s - %s", widget.CommonSettings().Title, widget.settings.langs)

	if widget.err != nil {
		return title, widget.err.Error(), true
	}

	if len(widget.repos) == 0 {
		return title, "No results to display", false
	}

	var str string
	for idx, repo := range widget.repos {
		var row []string
		rowNumber := fmt.Sprintf("[%s]%2d.", widget.RowColor(idx), idx+1)
		row = append(row, rowNumber)
		switch widget.showType {
		case ShowName:
			row = appendText(row, widget.settings.colors.repo, repo.Name)
		case ShowNameDesc:
			row = appendText(row, widget.settings.colors.repo, repo.Name)
			row = appendText(row, widget.RowColor(idx), repo.Description)
		case ShowNameDescLang:
			if repo.Language != "" {
				row = appendText(row, widget.settings.colors.lang, repo.Language)
			}
			row = appendText(row, widget.settings.colors.repo, repo.Name)
			row = appendText(row, widget.RowColor(idx), repo.Description)
		case ShowNameDescLangStars:
			if repo.Language != "" {
				row = appendText(row, widget.settings.colors.lang, repo.Language)
			}
			row = appendText(row, widget.settings.colors.repo, repo.Name)
			if repo.StarsToday != -1 {
				row = append(row, fmt.Sprintf(
					"[%s](⭐️%s⬆️%s)",
					widget.settings.colors.stars,
					formatThousands(repo.Stars),
					formatThousands(repo.StarsToday),
				))
			} else {
				row = append(row, fmt.Sprintf(
					"[%s](⭐️%s)",
					widget.settings.colors.stars,
					formatThousands(repo.Stars),
				))
			}
			row = appendText(row, widget.RowColor(idx), repo.Description)
		}

		str += utils.HighlightableHelper(widget.View, strings.Join(row, " "), idx, len(repo.Name))
	}

	return title, str, false
}

func appendText(row []string, color string, value string) []string {
	return append(row, fmt.Sprintf(
		`[%s]%s`,
		color,
		value,
	))
}

func formatThousands(value int) string {
	if value < 1000 {
		return strconv.Itoa(value)
	}
	floatValue := float64(value) / 1000
	roundedValue := math.Round(floatValue*10) / 10
	formattedValue := strconv.FormatFloat(roundedValue, 'f', 1, 64) + "k"
	return formattedValue
}

func (widget *Widget) openRepo() {
	story := widget.selectedRepo()
	if story != nil {
		utils.OpenFile(story.Url)
	}
}

func (widget *Widget) selectedRepo() *Repo {
	var repo *Repo

	sel := widget.GetSelected()
	if sel >= 0 && widget.repos != nil && sel < len(widget.repos) {
		repo = &widget.repos[sel]
	}

	return repo
}

func rotateShowType(showtype ShowType) ShowType {
	returnValue := ShowName
	switch showtype {
	case ShowName:
		returnValue = ShowNameDesc
	case ShowNameDesc:
		returnValue = ShowNameDescLang
	case ShowNameDescLang:
		returnValue = ShowNameDescLangStars
	}
	return returnValue
}

func (widget *Widget) toggleDisplayText() {
	widget.showType = rotateShowType(widget.showType)
	widget.Render()
}

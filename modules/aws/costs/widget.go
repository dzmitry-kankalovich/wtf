package awscosts

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
	"strings"
)

const (
	labelWidth = 14
)

type Widget struct {
	view.ScrollableWidget

	settings *Settings
	err      error
	client   *costexplorer.Client
}

func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := &Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.Common),
		settings:         settings,
	}

	cfg, err := getConfig(settings)
	if err != nil {
		widget.err = err
	}
	widget.client = costexplorer.NewFromConfig(*cfg)

	widget.SetRenderFunction(widget.Render)
	widget.initializeKeyboardControls()

	return widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	widget.Render()
}

func (widget *Widget) Render() {
	widget.Redraw(widget.content)
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() (string, string, bool) {
	title := fmt.Sprintf("AWS Costs - %s", widget.settings.alias)
	contents, err := getMonthlyCosts(widget.client)
	if err != nil {
		widget.err = err
		return "", "", false
	}

	if widget.settings.topN != -1 {
		topN, err := getTopN(widget.client, widget.settings.topN)
		if err != nil {
			widget.err = err
			return "", "", false
		}
		contents += fmt.Sprintf("\n--- Top %d services by cost ---", widget.settings.topN)
		for idx, service := range topN {
			var row []string
			rowNumber := fmt.Sprintf("[%s]%2d.", widget.RowColor(idx), idx+1)
			row = append(row, rowNumber)
			row = appendText(row, service)

			contents += "\n" + strings.Join(row, " ")
		}
	}

	return title, contents, false
}

func getMonthlyCosts(client *costexplorer.Client) (string, error) {
	mUnit, mAmount, err := getMTD(client)
	if err != nil {
		return "", err
	}
	fUnit, fAmount, err := getForecast(client)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%-*s %s %s\n%-*s %s %s",
		labelWidth, "Month-To-Date:", mUnit, mAmount,
		labelWidth, "Forecast:", fUnit, fAmount,
	), nil
}

func appendText(row []string, cost ServiceCost) []string {
	return append(row, fmt.Sprintf(
		`%-*s %s %s`,
		labelWidth,
		cost.Name,
		cost.Unit,
		cost.Amount,
	))
}

func (widget *Widget) openConsole() {
	utils.OpenFile("https://us-east-1.console.aws.amazon.com/billing/home#/")
}

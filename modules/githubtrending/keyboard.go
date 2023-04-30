package githubtrending

import "github.com/gdamore/tcell/v2"

func (widget *Widget) initializeKeyboardControls() {
	widget.InitializeHelpTextKeyboardControl(widget.ShowHelp)
	widget.InitializeRefreshKeyboardControl(widget.Refresh)

	widget.SetKeyboardChar("j", widget.Next, "Select next item")
	widget.SetKeyboardChar("k", widget.Prev, "Select previous item")
	widget.SetKeyboardChar("o", widget.openRepo, "Open repository in browser")
	widget.SetKeyboardChar("t", widget.toggleDisplayText, "Toggle display between name, "+
		"name+description, lang+name+description and lang+name+stars+description")

	widget.SetKeyboardKey(tcell.KeyDown, widget.Next, "Select next item")
	widget.SetKeyboardKey(tcell.KeyUp, widget.Prev, "Select previous item")
	widget.SetKeyboardKey(tcell.KeyEnter, widget.openRepo, "Open repository in browser")
	widget.SetKeyboardKey(tcell.KeyEsc, widget.Unselect, "Clear selection")
}

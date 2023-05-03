package awscosts

import "github.com/gdamore/tcell/v2"

func (widget *Widget) initializeKeyboardControls() {
	widget.InitializeHelpTextKeyboardControl(widget.ShowHelp)
	widget.InitializeRefreshKeyboardControl(widget.Refresh)

	widget.SetKeyboardChar("o", widget.openConsole, "Open AWS Console in browser")
	widget.SetKeyboardKey(tcell.KeyEnter, widget.openConsole, "Open AWS Console in browser")
}

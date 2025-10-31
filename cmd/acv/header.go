package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Header struct {
	grid *tview.Grid
}

func NewHeader(configServers []string) *Header {

	acDropdown = tview.NewDropDown().
		SetFocusedStyle(
			tcell.Style{}.
				Background(tcell.ColorBlue).
				Foreground(tcell.ColorBlack),
		).
		SetFieldStyle(
			tcell.Style{}.
				Background(tcell.ColorBlack).
				Foreground(tcell.ColorAntiqueWhite),
		).
		SetLabel("Server: ")

	acDropdown.
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				app.SetFocus(getKeysManager().keys)
				return nil
			}
			return event
		})
	acDropdown.SetBorder(true)

	for _, vaultUri := range configServers {
		acDropdown.AddOption(vaultUri, func() {
			connect(vaultUri)
			fetchSettings("*")
			updateKeysList()
			app.SetFocus(getKeysManager().keys)
		})
	}

	menu := NewShortcutMenu()

	logo := tview.NewTextArea().
		SetText(VANITY_LOGO, false).
		SetTextStyle(UIStyles.VanityLogoStyle)

	grid := tview.NewGrid().
		SetRows(3, 0).
		SetColumns(0, 23).
		AddItem(acDropdown, 0, 0, 1, 1, 0, 0, false).
		AddItem(menu.GetPrimitive(), 1, 0, 1, 1, 0, 0, false).
		AddItem(logo, 0, 1, 2, 1, 0, 0, false)
	grid.SetBackgroundColor(tcell.ColorBlack)

	return &Header{
		grid: grid,
	}
}

func (h *Header) GetPrimitive() tview.Primitive {
	return h.grid
}

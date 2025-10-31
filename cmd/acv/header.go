package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Header struct {
	grid       *tview.Grid
	escapeFunc func()
	acDropdown *tview.DropDown
}

func NewHeader(
	configServers []string,
	escapeFunc func(),
	serverSelectedFunc func(string),
) *Header {

	header := Header{
		escapeFunc: escapeFunc,
	}

	header.acDropdown = tview.NewDropDown().
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

	header.acDropdown.
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				escapeFunc()
				return nil
			}
			return event
		})
	header.acDropdown.SetBorder(true)

	for _, vaultUri := range configServers {
		header.acDropdown.AddOption(vaultUri, func() {
			serverSelectedFunc(vaultUri)
		})
	}

	menu := NewShortcutMenu()

	logo := tview.NewTextArea().
		SetText(VANITY_LOGO, false).
		SetTextStyle(UIStyles.VanityLogoStyle)

	header.grid = tview.NewGrid().
		SetRows(3, 0).
		SetColumns(0, 23).
		AddItem(header.acDropdown, 0, 0, 1, 1, 0, 0, false).
		AddItem(menu.GetPrimitive(), 1, 0, 1, 1, 0, 0, false).
		AddItem(logo, 0, 1, 2, 1, 0, 0, false)
	header.grid.SetBackgroundColor(tcell.ColorBlack)

	return &header
}

func (h *Header) GetPrimitive() tview.Primitive {
	return h.grid
}

func (h *Header) SelectFirstServer() {
	h.acDropdown.SetCurrentOption(0)
}

package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type KeysManager struct {
	box             *tview.Grid
	keys            *tview.Table
	keySelectedFunc func(string)
}

func NewKeysManager(
	keySelectedFunc func(string),
) KeysManager {
	box := tview.NewGrid()
	box.SetBorder(true)

	box.SetFocusFunc(func() {
		box.SetBorderColor(tcell.ColorBlue)
	}).SetBlurFunc(func() {
		box.SetBorderColor(tcell.ColorWhite)
	})

	keys := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		Select(0, 0)

	keys.SetFocusFunc(func() {
		style := tcell.Style{}.Foreground(tcell.ColorBlue).Bold(true).Background(tcell.ColorBlack)
		for row := range keys.GetRowCount() {
			for col := range keys.GetColumnCount() {
				keys.GetCell(row, col).SetStyle(style)
			}
		}
	}).SetBlurFunc(func() {
		style := tcell.Style{}.Foreground(tcell.ColorAntiqueWhite).Bold(false).Background(tcell.ColorBlack)
		for row := range keys.GetRowCount() {
			for col := range keys.GetColumnCount() {
				keys.GetCell(row, col).SetStyle(style)
			}
		}
	})

	box.AddItem(keys, 0, 0, 1, 1, 0, 0, false)

	manager := KeysManager{
		box:             box,
		keys:            keys,
		keySelectedFunc: keySelectedFunc,
	}

	keys.SetSelectedFunc(manager.settingSelected)

	keys.SetBorderPadding(1, 1, 1, 1)

	return manager
}

func (km KeysManager) updateKeys(settings []string) {
	km.keys.Clear()
	for row, setting := range settings {
		km.keys.SetCell(row, 0, tview.NewTableCell(setting))
	}
}

func (km KeysManager) settingSelected(row int, col int) {
	settingCell := km.keys.GetCell(row, col)
	if settingCell == nil {
		return
	}

	km.keySelectedFunc(settingCell.Text)
}

func (km KeysManager) SetTitle(title string) {
	km.box.SetTitle(title)
}

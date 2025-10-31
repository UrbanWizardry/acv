package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type KeysManager struct {
	grid                 *tview.Grid
	keys                 *tview.Table
	settingSearchManager *SearchManager
	keySelectedFunc      func(string)
}

func NewKeysManager(
	keySelectedFunc func(string),
) *KeysManager {
	manager := KeysManager{
		keySelectedFunc: keySelectedFunc,
	}

	manager.grid = tview.NewGrid()

	// Set things that chain as *tview.Box
	manager.grid.
		SetRows(0, 3).
		SetFocusFunc(func() {
			manager.grid.SetBorderColor(tcell.ColorBlue)
		}).SetBlurFunc(func() {
		manager.grid.SetBorderColor(tcell.ColorWhite)
	})

	manager.keys = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		Select(0, 0).
		SetSelectedFunc(manager.settingSelected)

	// Set things that chain as *tview.Box
	manager.keys.SetFocusFunc(func() {
		style := tcell.Style{}.Foreground(tcell.ColorBlue).Bold(true).Background(tcell.ColorBlack)
		for row := range manager.keys.GetRowCount() {
			for col := range manager.keys.GetColumnCount() {
				manager.keys.GetCell(row, col).SetStyle(style)
			}
		}
	}).SetBlurFunc(func() {
		style := tcell.Style{}.Foreground(tcell.ColorAntiqueWhite).Bold(false).Background(tcell.ColorBlack)
		for row := range manager.keys.GetRowCount() {
			for col := range manager.keys.GetColumnCount() {
				manager.keys.GetCell(row, col).SetStyle(style)
			}
		}
	}).SetBorderPadding(1, 1, 1, 1)

	keysBox := tview.NewGrid()
	keysBox.SetBorder(true)
	keysBox.AddItem(manager.keys, 0, 0, 1, 1, 0, 0, false)

	manager.grid.AddItem(keysBox, 0, 0, 1, 1, 0, 0, false)

	// Setting Search Bar
	manager.settingSearchManager = NewSearchManager(
		func(p tview.Primitive) {
			app.SetFocus(p)
		},
		func(s string) {
			findSettings(s)
			updateKeysList()
			app.SetFocus(manager.keys)
			manager.settingSearchManager.setSearchType(NoSearch)
		},
	)

	manager.grid.AddItem(manager.settingSearchManager.GetPrimitive(), 1, 0, 1, 1, 0, 0, false)

	return &manager
}

func (km *KeysManager) GetPrimitive() tview.Primitive {
	return km.grid
}

func (km *KeysManager) updateKeys(settings []string) {
	km.keys.Clear()
	for row, setting := range settings {
		km.keys.SetCell(row, 0, tview.NewTableCell(setting))
	}
}

func (km *KeysManager) settingSelected(row int, col int) {
	settingCell := km.keys.GetCell(row, col)
	if settingCell == nil {
		return
	}

	km.keySelectedFunc(settingCell.Text)
}

func (km *KeysManager) SetTitle(title string) {
	km.grid.SetTitle(title)
}

package main

import (
	"fmt"
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ShortcutMenu struct {
	grid *tview.Grid
}

func NewShortcutMenu() ShortcutMenu {

	menu := ShortcutMenu{
		grid: tview.NewGrid(),
	}

	cols := []map[rune]string{
		map[rune]string{
			's': "Change config server",
			'q': "Quit ACV",
		},
		map[rune]string{
			'/': "Search (keys or value)",
			'r': "Reload keys list",
		},
		map[rune]string{
			'j': "Toggle JSON prettyprint",
			'd': "Toggle diff mode",
			'c': "Copy selected value",
		},
	}

	// Use two more rows than needed to create padding.
	maxRows := slices.Max(arraymap(cols, func(c map[rune]string) int { return len(c) }))
	menu.grid.SetRows(func(len int) []int {
		r := []int{}
		for range len {
			r = append(r, 1)
		}
		return r
	}(maxRows + 2)...)

	menu.grid.SetBackgroundColor(tcell.ColorBlack)

	for c, col := range cols {
		menu.grid.AddItem(tview.NewBox(), 0, c, 1, 1, 0, 0, false)

		row := 1
		for key, item := range col {
			menu.grid.AddItem(tview.NewInputField().SetText(item).SetLabel(fmt.Sprintf("<%c> ", key)).SetDisabled(true), row, c, 1, 1, 0, 0, false)
			row++
		}

		for i := range 4 - len(col) {
			menu.grid.AddItem(tview.NewBox(), len(col)+i+1, c, 1, 1, 0, 0, false)
		}
	}

	return menu
}

func (sm *ShortcutMenu) GetPrimitive() tview.Primitive {
	return sm.grid
}

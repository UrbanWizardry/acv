package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var searchLabels = map[SearchType]string{
	NoSearch:     "Search: ",
	StringSearch: "Search: ",
}

type SearchManager struct {
	searchBox         *tview.InputField
	searchType        SearchType
	setFocusFunc      func(tview.Primitive)
	searchChangedFunc func(string)
}

func NewSearchManager(
	setFocusFunc func(tview.Primitive),
	searchChangedFunc func(string),
) SearchManager {
	// Setting search box
	searchBox := tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabel(searchLabels[NoSearch])

	searchBox.SetBorder(true).
		SetFocusFunc(func() {
			searchBox.
				SetFieldStyle(tcell.Style{}.
					Background(tcell.ColorBlack).
					Foreground(tcell.ColorAntiqueWhite),
				).
				SetBorderColor(tcell.ColorBlue)
		}).
		SetBlurFunc(func() {
			searchBox.
				SetFieldStyle(tcell.Style{}.
					Background(tcell.ColorBlack).
					Foreground(tcell.ColorAntiqueWhite),
				).
				SetBorderColor(tcell.ColorWhite)
		})

	manager := SearchManager{
		searchBox:         searchBox,
		searchType:        NoSearch,
		setFocusFunc:      setFocusFunc,
		searchChangedFunc: searchChangedFunc,
	}

	searchBox.SetInputCapture(manager.onInput)

	return manager
}

func (sm *SearchManager) GetPrimitive() tview.Primitive {
	return sm.searchBox
}

func (sm *SearchManager) setSearchType(st SearchType) {
	sm.searchType = st
	sm.searchBox.SetLabel(searchLabels[sm.searchType])
}

func (sm *SearchManager) setSearching(st SearchType) {
	sm.setSearchType(st)
	sm.setFocusFunc(sm.searchBox)
}

func (sm *SearchManager) exitSearching() {
	sm.setSearchType(NoSearch)
}

func (sm *SearchManager) onInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		// Leave searching, clear search
		//fetchSettings("*")
		sm.exitSearching()
		sm.searchChangedFunc("")
		return nil

	case tcell.KeyEnter:
		// Leave searching, apply search
		//fetchSettings(fmt.Sprintf("%s*", sm.searchBox.GetText()))
		sm.exitSearching()
		sm.searchChangedFunc(sm.searchBox.GetText())
		return nil
	}

	return event
}

func (sm *SearchManager) Reset() {
	sm.searchBox.SetText("")
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	//"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/kylelemons/godebug/diff"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var valueTitles = map[RenderType]string{
	Plain: "Formatting: Plain",
	Json:  "Formatting: JSON",
}

type ValuesManager struct {
	// Regular display of which setting/versions is selected
	primaryRevisionSelector *ValuesRevisionSelector
	diffRevisionSelector    *ValuesRevisionSelector

	// Revision value display
	valueTextView      *tview.TextView
	valueSearchManager *SearchManager

	// UI Layout
	grid *tview.Grid

	// Internal state

	renderType   RenderType
	setFocusFunc func(tview.Primitive)
}

var _ UIComponent = (*ValuesManager)(nil)

func NewValuesManager(
	escapeFunc func(),
	setFocusFunc func(tview.Primitive),
) *ValuesManager {
	manager := &ValuesManager{
		setFocusFunc: setFocusFunc,
	}

	// Primary revision selector, for setting view or diff view left value
	primaryRevisionSelector := NewValuesRevisionSelector(
		func() {
			escapeFunc()
		},
		func(value string) {
			manager.updateValueBasedOnView()
		},
		setFocusFunc,
		Left,
	)
	manager.primaryRevisionSelector = primaryRevisionSelector

	// Primary revision selector, for setting view or diff view left value
	diffRightRevisionSelector := NewValuesRevisionSelector(
		func() {
			primaryRevisionSelector.focusRevisionDropdown()
		},
		func(value string) {
			// This can only happen in Diff mode, so just get on with it
			// diff values already takes care of fomatting, don't do it here
			manager.updateValueBasedOnView()
			//manager.updateValue(manager.diffValues())
		},
		setFocusFunc,
		Right,
	)
	manager.diffRevisionSelector = diffRightRevisionSelector

	// Config view for displaying setting value of diff output
	configValue := tview.NewTextView().SetRegions(true).SetDynamicColors(true)
	configValue.SetDisabled(true)
	configValue.
		SetBorderPadding(1, 1, 1, 1).
		SetBorder(true).
		SetFocusFunc(func() {
			configValue.SetBorderColor(tcell.ColorBlue)
		}).
		SetBlurFunc(func() {
			configValue.SetBorderColor(tcell.ColorWhite)
		})

	configValue.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			if manager.primaryRevisionSelector.viewMode == Standard {
				primaryRevisionSelector.focusRevisionDropdown()
			} else {
				diffRightRevisionSelector.focusRevisionDropdown()
			}
			return nil
		}
		return event
	})

	manager.valueTextView = configValue

	// Value Text Search Bar
	manager.valueSearchManager = NewSearchManager(
		func(p tview.Primitive) {
			setFocusFunc(p)
		},
		func(s string) {
			if s == "" {
				manager.updateValueBasedOnView()
			} else {
				manager.updateValue(strings.Replace(manager.getValueBasedOnView(), s, fmt.Sprintf("[\"search\"]%s[\"\"]", s), -1))
			}

			configValue.Highlight("search")
			configValue.ScrollToHighlight()
			manager.valueSearchManager.setSearchType(NoSearch)
		},
	)

	// Layout Grid
	grid := tview.NewGrid()
	manager.grid = grid
	manager.layoutStandard()

	return manager
}

func (vm *ValuesManager) GetPrimitive() tview.Primitive {
	return vm.grid
}

func (vm *ValuesManager) SetDisplayMode(mode ValueDisplayMode) {

	vm.primaryRevisionSelector.setDisplayMode(mode)
	vm.diffRevisionSelector.setDisplayMode(mode)

	switch mode {
	case Standard:
		vm.layoutStandard()
		// Restore a standard view
		vm.updateValue(vm.formatValue(vm.primaryRevisionSelector.GetCurrentValue()))
		vm.diffRevisionSelector.Clear()
	case Diff:
		vm.layoutDiff()
	}
}

func (vm *ValuesManager) layoutStandard() {
	vm.grid.Clear()
	vm.grid.
		SetRows(3, 0, 3).
		AddItem(vm.primaryRevisionSelector.GetPrimitive(), 0, 0, 1, 1, 0, 0, false).
		AddItem(vm.valueTextView, 1, 0, 1, 1, 0, 0, false).
		AddItem(vm.valueSearchManager.GetPrimitive(), 2, 0, 1, 1, 0, 0, false)
}

func (vm *ValuesManager) layoutDiff() {
	vm.grid.Clear()
	vm.grid.
		SetRows(3, 3, 0, 3).
		AddItem(vm.primaryRevisionSelector.GetPrimitive(), 0, 0, 1, 1, 0, 0, false).
		AddItem(vm.diffRevisionSelector.GetPrimitive(), 1, 0, 1, 1, 0, 0, false).
		AddItem(vm.valueTextView, 2, 0, 1, 1, 0, 0, false).
		AddItem(vm.valueSearchManager.GetPrimitive(), 3, 0, 1, 1, 0, 0, false)
}

func (vm *ValuesManager) setTextViewTitle() {
	vm.valueTextView.SetTitle(valueTitles[vm.renderType])
}

func (vm *ValuesManager) setRenderType(t RenderType) {
	vm.renderType = t
	vm.setTextViewTitle()
}

func (vm *ValuesManager) reset() {
	vm.primaryRevisionSelector.setRevisions("", []azappconfig.Setting{})
	vm.valueTextView.SetText("")
}

func (vm *ValuesManager) setPrimaryRevisions(settingName string, revisions []azappconfig.Setting) {
	vm.primaryRevisionSelector.setRevisions(settingName, revisions)
}

func (vm *ValuesManager) setDiffRightRevisions(settingName string, revisions []azappconfig.Setting) {
	vm.diffRevisionSelector.setRevisions(settingName, revisions)
}

func (vm *ValuesManager) updateValueToPrimaryRevision() {
	if len(vm.primaryRevisionSelector.revisions) > 0 {
		vm.updateValue(vm.formatValue(vm.primaryRevisionSelector.GetCurrentValue()))
	}
}

func (vm *ValuesManager) updateValueBasedOnView() {
	vm.updateValue(vm.getValueBasedOnView())
}

func (vm *ValuesManager) getValueBasedOnView() string {
	// If Standard mode, format before return the value.
	// If Diff mode, do the diff thing (does formatting for you)
	if vm.primaryRevisionSelector.viewMode == Standard {
		return vm.formatValue(vm.primaryRevisionSelector.GetCurrentValue())
	} else {
		return vm.diffValues()
	}
}

func (vm *ValuesManager) formatValue(value string) string {
	var printValue string
	switch vm.renderType {
	case Plain:
		printValue = value
	case Json:
		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, []byte(value), "", "   ")
		if error != nil {
			// Can't format as JSON, default to plain
			printValue = value
		} else {
			printValue = prettyJSON.String()
		}
	}

	return printValue
}

// updateValue sets the text view to hold this exact string with no more fortmatting,
// then fiddles with focus and UI aspects
func (vm *ValuesManager) updateValue(value string) {
	vm.setValue(value)
	vm.setFocusFunc(vm.valueTextView)
	vm.setTextViewTitle()
}

func (vm *ValuesManager) setValue(value string) {
	vm.valueTextView.SetText(value)
}

func (vm *ValuesManager) diffValues() string {
	s := diff.Diff(
		vm.formatValue(vm.primaryRevisionSelector.GetCurrentValue()),
		vm.formatValue(vm.diffRevisionSelector.GetCurrentValue()),
	)
	lines := strings.Split(s, "\n")

	formatlines := arraymap(lines, func(s string) string {
		if strings.HasPrefix(s, "-") {
			return fmt.Sprintf("[red]%s[white]", s)
		} else if strings.HasPrefix(s, "+") {
			return fmt.Sprintf("[green]%s[white]", s)
		} else {
			return s
		}
	})

	return strings.Join(formatlines, "\n")
}

// UTILITY

func sortRevisionsNewestFirst(versions []azappconfig.Setting) []azappconfig.Setting {
	sort.Slice(versions, func(i, j int) bool {
		// We return the inverse of "less", because we want descending order
		return !versions[i].LastModified.Before(*versions[j].LastModified)
	})

	return versions
}

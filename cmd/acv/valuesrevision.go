package main

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ValuesRevisionSelector struct {
	// UI Layout
	revisionsGrid         *tview.Grid
	revisionsSettingLabel *tview.InputField
	revisionsDropDown     *tview.DropDown

	// Events and Callbacks
	revisionChangedFunc func(string)
	escapeFunc          func()
	setFocusFunc        func(tview.Primitive)

	// Internal State
	revisions  []azappconfig.Setting
	viewMode   ValueDisplayMode
	diffSource DiffSource
}

func NewValuesRevisionSelector(
	escapeFunc func(),
	revisionChangedFunc func(string),
	setFocusFunc func(tview.Primitive),
	diffSource DiffSource,
) *ValuesRevisionSelector {

	selector := &ValuesRevisionSelector{
		revisionChangedFunc: revisionChangedFunc,
		escapeFunc:          escapeFunc,
		setFocusFunc:        setFocusFunc,
		diffSource:          diffSource,
	}

	revisionsSettingLabel := tview.NewInputField()
	revisionsSettingLabel.SetDisabled(true)
	revisionsSettingLabel.SetBorder(true)
	selector.revisionsSettingLabel = revisionsSettingLabel

	revisionsDropDown := tview.NewDropDown().SetLabel("Rev: ").SetFieldBackgroundColor(tcell.ColorBlack)
	revisionsDropDown.
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				escapeFunc()
				return nil
			}
			return event
		}).
		SetBorder(true).
		SetFocusFunc(selector.setFocusStyle).
		SetBlurFunc(selector.setBlurStyle)
	selector.revisionsDropDown = revisionsDropDown

	revisionsGrid := tview.NewGrid().SetColumns(-2, -1)
	revisionsGrid.AddItem(revisionsSettingLabel, 0, 0, 1, 1, 0, 0, false)
	revisionsGrid.AddItem(revisionsDropDown, 0, 1, 1, 1, 0, 0, false)
	selector.revisionsGrid = revisionsGrid

	return selector
}

func (vrs *ValuesRevisionSelector) GetPrimitive() tview.Primitive {
	return vrs.revisionsGrid
}

func (vrs *ValuesRevisionSelector) setRevisions(settingName string, revisions []azappconfig.Setting) {
	vrs.revisionsSettingLabel.SetText(settingName)
	vrs.revisions = revisions
	vrs.revisionsDropDown.SetOptions([]string{}, nil)
	if len(revisions) > 0 {
		sortedRevisions := sortRevisionsNewestFirst(vrs.revisions)
		revisionOptions := arraymap(
			sortedRevisions,
			func(r azappconfig.Setting) string { return r.LastModified.Format(time.RFC822) },
		)

		vrs.revisionsDropDown.SetOptions(
			revisionOptions,
			vrs.revisionSelected,
		)

		vrs.revisionsDropDown.SetCurrentOption(0)
	}
}

func (vrs *ValuesRevisionSelector) revisionSelected(text string, index int) {
	vrs.revisionChangedFunc(*vrs.revisions[index].Value)
}

func (vrs *ValuesRevisionSelector) setDisplayMode(mode ValueDisplayMode) {
	vrs.viewMode = mode
	if vrs.revisionsDropDown.HasFocus() {
		vrs.setFocusStyle()
	} else {
		vrs.setBlurStyle()
	}
}

// Focus control

// focusRevisionDropdown is used externally to focus the dropdown in this selector
func (vrs *ValuesRevisionSelector) focusRevisionDropdown() {
	vrs.setFocusFunc(vrs.revisionsDropDown)
}

func (vrs *ValuesRevisionSelector) setFocusStyle() {
	vrs.revisionsSettingLabel.SetBorderStyle(UIStyles.RevisionSelectorBorderFocus)
	vrs.revisionsDropDown.SetBorderStyle(UIStyles.RevisionSelectorBorderFocus)
}

func (vrs *ValuesRevisionSelector) setBlurStyle() {

	// Not focused, so choose a blur style
	if vrs.viewMode == Diff {
		if vrs.diffSource == Left {
			vrs.revisionsSettingLabel.SetBorderStyle(UIStyles.RevisionSelectorBorderDiffLeft)
			vrs.revisionsDropDown.SetBorderStyle(UIStyles.RevisionSelectorBorderDiffLeft)
		} else {
			vrs.revisionsSettingLabel.SetBorderStyle(UIStyles.RevisionSelectorBorderDiffRight)
			vrs.revisionsDropDown.SetBorderStyle(UIStyles.RevisionSelectorBorderDiffRight)
		}
	} else {
		vrs.revisionsSettingLabel.SetBorderStyle(UIStyles.RevisionSelectorBorderBlur)
		vrs.revisionsDropDown.SetBorderStyle(UIStyles.RevisionSelectorBorderBlur)
	}

}

func (vrs *ValuesRevisionSelector) GetCurrentValue() string {
	index, _ := vrs.revisionsDropDown.GetCurrentOption()
	return *vrs.revisions[index].Value
}

func (vrs *ValuesRevisionSelector) Clear() {
	vrs.revisionsSettingLabel.SetText("")
	vrs.revisionsDropDown.SetOptions([]string{}, nil)
}

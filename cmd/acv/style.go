package main

import "github.com/gdamore/tcell/v2"

type UIStylesDef struct {
	// Vanity
	VanityLogoStyle tcell.Style

	// Dropdowns
	DropdownFocus tcell.Style
	DropdownBlur  tcell.Style

	// Table cells i.e. for lists
	TableCellFocus tcell.Style
	TableCellBlur  tcell.Style

	// Revision Selectors
	RevisionSelectorBorderBlur      tcell.Style
	RevisionSelectorBorderFocus     tcell.Style
	RevisionSelectorBorderDiffLeft  tcell.Style
	RevisionSelectorBorderDiffRight tcell.Style
}

var UIStyles UIStylesDef = UIStylesDef{
	VanityLogoStyle: tcell.Style{}.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorBlue),

	DropdownFocus: tcell.Style{}.
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorBlack),

	DropdownBlur: tcell.Style{}.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorAntiqueWhite),

	TableCellFocus: tcell.Style{}.
		Foreground(tcell.ColorBlue).
		Bold(true).
		Background(tcell.ColorBlack),

	TableCellBlur: tcell.Style{}.
		Foreground(tcell.ColorAntiqueWhite).
		Bold(false).
		Background(tcell.ColorBlack),

	RevisionSelectorBorderBlur: tcell.Style{}.
		Foreground(tcell.ColorAntiqueWhite).
		Background(tcell.ColorBlack),

	RevisionSelectorBorderFocus: tcell.Style{}.
		Foreground(tcell.ColorBlue).
		Background(tcell.ColorBlack),

	RevisionSelectorBorderDiffLeft: tcell.Style{}.
		Foreground(tcell.ColorRed).
		Background(tcell.ColorBlack),

	RevisionSelectorBorderDiffRight: tcell.Style{}.
		Foreground(tcell.ColorGreen).
		Background(tcell.ColorBlack),
}

package main

import "github.com/rivo/tview"

type SearchType int

const (
	NoSearch SearchType = iota
	StringSearch
)

type RenderType int

const (
	Plain RenderType = iota
	Json
)

type ValueDisplayMode int

const (
	Standard ValueDisplayMode = iota
	Diff
)

type DiffSource int

const (
	Left DiffSource = iota
	Right
)

type UIComponent interface {
	GetPrimitive() tview.Primitive
}

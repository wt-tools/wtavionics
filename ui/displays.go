package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type basicDisplay struct {
	V      string
	title  string
	valw   material.LabelStyle
	theme  *material.Theme
	height int
}

func newBasicDisplay(theme *material.Theme, title string, height int) *basicDisplay {
	return &basicDisplay{
		title:  title,
		valw:   material.Label(theme, unit.Sp(height), "0"),
		theme:  theme,
		height: height,
	}
}

func (b *basicDisplay) Layout(gtx C) D {
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(b.theme, unit.Sp(30), b.title)
				lbl.Color = color.NRGBA{}
				if b.V == "" || b.V == "0" {
					lbl.Color = color.NRGBA{100, 100, 100, 96}
				}
				lbl.Alignment = text.Middle
				return lbl.Layout(gtx)
			}),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				b.valw.Text = b.V
				b.valw.Alignment = text.Middle
				return b.valw.Layout(gtx)
			},
		),
	)
}

package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func newCompassDisplay(theme *material.Theme) *compassDisplay {
	return &compassDisplay{theme: theme}
}

type compassDisplay struct {
	//	director[int] TODO доработать director для left-right not only up-down
	v     int
	label string
	theme *material.Theme
	color color.NRGBA
}

func (c *compassDisplay) Set(d int) {
	c.v = d
	switch d {
	case 0:
		c.label = "N"
	case 90:
		c.label = "E"
	case 180:
		c.label = "S"
	case 270:
		c.label = "W"
	default:
		switch {
		case d > 0 && d < 90:
			c.label = fmt.Sprintf("N %d E", d)
		case d > 90 && d < 180:
			c.label = fmt.Sprintf("E %d S", d)
		case d > 180 && d < 270:
			c.label = fmt.Sprintf("S %d W", d)
		case d > 270 && d < 360:
			c.label = fmt.Sprintf("W %d N", d)
		}
	}
}

func (c *compassDisplay) Display(gtx C, visible bool) func(C) D {
	switch {
	case !visible:
		c.color = color.NRGBA{0, 0, 0, 0}
	case c.v == 0:
		c.color = color.NRGBA{100, 100, 100, 70}
	default:
		c.color = color.NRGBA{0, 0, 0, 255}
	}
	return func(gtx C) D {
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceEnd,
		}.Layout(gtx,
			layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(c.theme, unit.Sp(60), c.label)
					lbl.Color = c.color
					lbl.Alignment = text.Middle
					return lbl.Layout(gtx)
				}),
		)
	}
}

func newIntBasicDisplay(theme *material.Theme, title string, height int) *intBasicDisplay {
	d := basicDisplay{
		title:  title,
		label:  title,
		valw:   material.Label(theme, unit.Sp(height), "0"),
		theme:  theme,
		height: height,
	}
	return &intBasicDisplay{basicDisplay: d}
}

func newFloatBasicDisplay(theme *material.Theme, title string, height int) *floatBasicDisplay {
	d := basicDisplay{
		title:  title,
		label:  title,
		valw:   material.Label(theme, unit.Sp(height), "0"),
		theme:  theme,
		height: height,
	}
	return &floatBasicDisplay{basicDisplay: d}
}

type intBasicDisplay struct {
	basicDisplay
	director[int]
}

func (i *intBasicDisplay) Set(d int) {
	i.V = strconv.Itoa(d)
	i.label = i.title + " " + i.set(d)
}

type floatBasicDisplay struct {
	basicDisplay
	director[float64]
}

func (f *floatBasicDisplay) Set(d float64) {
	f.V = strconv.FormatFloat(d, 'f', precision, 64)
	f.label = f.title + " " + f.set(d)
}

type basicDisplay struct {
	V      string
	title  string
	label  string
	valw   material.LabelStyle
	theme  *material.Theme
	height int
	color  color.NRGBA
}

func (b *basicDisplay) Display(gtx C, visible bool) func(C) D {
	switch {
	case !visible:
		b.color = color.NRGBA{0, 0, 0, 0}
	case b.V == "" || b.V == "0" || b.V == "0.0":
		b.color = color.NRGBA{100, 100, 100, 70}
	default:
		b.color = color.NRGBA{0, 0, 0, 255}
	}
	return func(gtx C) D {
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceEnd,
		}.Layout(gtx,
			layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(b.theme, unit.Sp(30), b.label)
					lbl.Color = b.color
					lbl.Alignment = text.Middle
					return lbl.Layout(gtx)
				}),
			layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					b.valw.Text = b.V
					b.valw.Color = b.color
					b.valw.Alignment = text.Middle
					return b.valw.Layout(gtx)
				},
			),
		)
	}
}

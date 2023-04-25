package ui

import (
	"context"
	"fmt"
	"strconv"

	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/wt-tools/wtscope/input/indicators"
	"github.com/wt-tools/wtscope/input/state"
)

type avionicsDisplays struct {
	// Updatable elements:
	speed        material.LabelStyle
	altH         material.LabelStyle
	oilTemp      material.LabelStyle
	headTemp     material.LabelStyle
	fuel         material.LabelStyle
	flaps        material.LabelStyle
	craft        material.ButtonStyle
	btnClickArea widget.Clickable
}

func (g *gui) UpdateAvionics(ctx context.Context, states *state.Service, inds *indicators.Service) {
	l := g.log.New()
	g.speed = material.Label(g.th, unit.Sp(50), "0")
	g.altH = material.Label(g.th, unit.Sp(50), "0")
	g.oilTemp = material.Label(g.th, unit.Sp(50), "0")
	g.headTemp = material.Label(g.th, unit.Sp(50), "0")
	g.fuel = material.Label(g.th, unit.Sp(50), "0")
	g.flaps = material.Label(g.th, unit.Sp(50), "off")
	g.craft = material.Button(g.th, &g.btnClickArea, "aircraft")
	go func() {
		for {
			select {
			case data := <-states.Messages:
				l.Log("state", data)
			case data := <-inds.Messages:
				g.speed.Text = strconv.FormatFloat(data.Speed, 'f', 3, 64)
				g.altH.Text = strconv.FormatFloat(data.AltitudeHour, 'f', 3, 64)
				g.oilTemp.Text = strconv.FormatFloat(data.OilTemperature, 'f', 3, 64)
				g.headTemp.Text = strconv.FormatFloat(data.HeadTemperature, 'f', 3, 64)
				g.fuel.Text = strconv.FormatFloat(data.Fuel, 'f', 3, 64)
				g.flaps.Text = strconv.FormatFloat(data.Flaps, 'f', 3, 64)
				g.craft.Text = data.Type
				g.w.Invalidate()
				l.Log("indicator", data)
			}
		}
	}()
}

func (g *gui) avionicsPanel() error {
	l := g.log.New()
	var ops op.Ops
	btn1 := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			return g.craft.Layout(gtx)
		},
	)

	btn2 := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(g.th, &g.btnClickArea, "hello world")

			return btn.Layout(gtx)
		},
	)

	rows := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}

	for e := range g.w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
			for _, event := range gtx.Events(g.w) {
				switch event := event.(type) {
				case key.Event:
					l.Log("event", "key pressed")
					if event.Name == key.NameEscape {
						return nil
					}
				}
			}
			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				layout.Rigid(
					layout.Spacer{Height: unit.Dp(250)}.Layout,
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(g.th, unit.Sp(30), "speed")
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						g.speed.Alignment = text.Middle
						return g.speed.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(g.th, unit.Sp(30), "altitude")
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						g.altH.Alignment = text.Middle
						return g.altH.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(g.th, unit.Sp(30), "oil temperature")
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						g.oilTemp.Alignment = text.Middle
						return g.oilTemp.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(g.th, unit.Sp(30), "head temperature")
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						g.headTemp.Alignment = text.Middle
						return g.headTemp.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(g.th, unit.Sp(30), "fuel")
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						g.fuel.Alignment = text.Middle
						return g.fuel.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(g.th, unit.Sp(30), "flaps")
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						g.flaps.Alignment = text.Middle
						return g.flaps.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return rows.Layout(gtx, btn1, btn2)
					},
				),
			)
			if g.btnClickArea.Clicked() {
				fmt.Println("button was clicked")
			}
			area.Pop()
			e.Frame(gtx.Ops)
		}
	}

	return nil
}

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
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/wt-tools/wtscope/input/indicators"
	"github.com/wt-tools/wtscope/input/state"
)

const noAircraft = "no aircraft selected"

type avionicsDisplays struct {
	ias          *basicDisplay
	altH         *basicDisplay
	oilTemp      *basicDisplay
	headTemp     *basicDisplay
	waterTemp    *basicDisplay
	fuel         *basicDisplay
	flaps        *basicDisplay
	throttle     *basicDisplay
	craft        material.ButtonStyle
	btnClickArea widget.Clickable
}

const precision = 1 // numbers after comma for floating values
func (g *gui) UpdateAvionics(ctx context.Context, states *state.Service, inds *indicators.Service) {
	l := g.log.New()
	g.ias = newBasicDisplay(g.th, "speed", 330)
	g.altH = newBasicDisplay(g.th, "altitude", 50)
	g.oilTemp = newBasicDisplay(g.th, "oil temperature", 90)
	g.waterTemp = newBasicDisplay(g.th, "water temperature", 90)
	g.headTemp = newBasicDisplay(g.th, "head temperature", 90)
	g.fuel = newBasicDisplay(g.th, "fuel", 50)
	g.flaps = newBasicDisplay(g.th, "flaps", 50)
	g.throttle = newBasicDisplay(g.th, "throttle", 50)
	g.craft = material.Button(g.th, &g.btnClickArea, noAircraft)
	go func() {
		for {
			select {
			case data := <-states.Messages:
				g.altH.V = strconv.Itoa(data.GetInt(state.HM))
				g.ias.V = strconv.Itoa(data.GetInt(state.IASKmH))
				g.throttle.V = strconv.Itoa(data.GetInt(state.Throttle1))
				g.w.Invalidate()
				l.Log("state", data)
			case data := <-inds.Messages:
				if data.OilTemperature >= 0 {
					g.oilTemp.V = strconv.FormatFloat(data.OilTemperature, 'f', precision, 64)
				}
				if data.HeadTemperature >= 0 {
					g.headTemp.V = strconv.FormatFloat(data.HeadTemperature, 'f', precision, 64)
				}
				if data.WaterTemperature >= 0 {
					g.waterTemp.V = strconv.FormatFloat(data.WaterTemperature, 'f', precision, 64)
				}
				g.fuel.V = strconv.FormatFloat(data.Fuel, 'f', precision, 64)
				g.flaps.V = strconv.FormatFloat(data.Flaps, 'f', precision, 64)
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
			btn := material.Button(g.th, &g.btnClickArea, "dev preview")

			return btn.Layout(gtx)
		},
	)

	rows := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}

	for e := range g.w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
			visible := !(g.craft.Text == noAircraft) || g.craft.Text == ""
			for _, event := range gtx.Events(g.w) {
				switch event := event.(type) {
				case key.Event:
					l.Log("exit", "by escape")
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
					layout.Spacer{Height: unit.Dp(25)}.Layout,
				),
				layout.Rigid(g.ias.Display(gtx, visible)),
				layout.Rigid(g.altH.Display(gtx, visible)),
				layout.Rigid(g.oilTemp.Display(gtx, visible)),
				layout.Rigid(g.waterTemp.Display(gtx, visible)),
				layout.Rigid(g.headTemp.Display(gtx, visible)),
				layout.Rigid(g.throttle.Display(gtx, visible)),
				layout.Rigid(g.flaps.Display(gtx, visible)),
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

// func (g *gui) basicIndicators(names ...string) []layout.FlexChild {
//	var children []layout.FlexChild
//	for _, name := range names {
//		children = append(children,
//			layout.Rigid(
//				func(gtx layout.Context) layout.Dimensions {
//					lbl := material.Label(g.th, unit.Sp(30), name)
//					lbl.Alignment = text.Middle
//					return lbl.Layout(gtx)
//				}),
//			layout.Rigid(
//				func(gtx layout.Context) layout.Dimensions {
//					g.ias.Alignment = text.Middle
//					return g.ias.Layout(gtx)
//				}),
//		)
//	}
//	return children
// }

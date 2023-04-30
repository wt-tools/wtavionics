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
	//	displays map[string]  material.LabelStyle
	ias          material.LabelStyle
	altH         material.LabelStyle
	oilTemp      material.LabelStyle
	headTemp     material.LabelStyle
	waterTemp    material.LabelStyle
	fuel         material.LabelStyle
	flaps        material.LabelStyle
	throttle     material.LabelStyle
	craft        material.ButtonStyle
	btnClickArea widget.Clickable
}

const precision = 1 // numbers after comma for floating values
func (g *gui) UpdateAvionics(ctx context.Context, states *state.Service, inds *indicators.Service) {
	l := g.log.New()
	//	g.displays = make(map[string]material.LabelStyle)
	g.ias = material.Label(g.th, unit.Sp(300), "0")
	g.altH = material.Label(g.th, unit.Sp(50), "0")
	g.oilTemp = material.Label(g.th, unit.Sp(70), "0")
	g.headTemp = material.Label(g.th, unit.Sp(70), "0")
	g.waterTemp = material.Label(g.th, unit.Sp(70), "0")
	g.fuel = material.Label(g.th, unit.Sp(50), "0")
	g.flaps = material.Label(g.th, unit.Sp(50), "off")
	g.throttle = material.Label(g.th, unit.Sp(50), "0")
	g.craft = material.Button(g.th, &g.btnClickArea, "aircraft")
	go func() {
		for {
			select {
			case data := <-states.Messages:
				g.altH.Text = strconv.Itoa(data.GetInt(state.HM))
				g.ias.Text = strconv.Itoa(data.GetInt(state.IASKmH))
				g.throttle.Text = strconv.Itoa(data.GetInt(state.Throttle1))
				g.w.Invalidate()
				l.Log("state", data)
			case data := <-inds.Messages:
				if data.OilTemperature < 0 {
					data.OilTemperature = 0
				}
				g.oilTemp.Text = strconv.FormatFloat(data.OilTemperature, 'f', precision, 64)
				if data.HeadTemperature < 0 {
					data.HeadTemperature = 0
				}
				g.headTemp.Text = strconv.FormatFloat(data.HeadTemperature, 'f', precision, 64)
				if data.WaterTemperature < 0 {
					data.WaterTemperature = 0
				}
				g.waterTemp.Text = strconv.FormatFloat(data.WaterTemperature, 'f', precision, 64)
				g.fuel.Text = strconv.FormatFloat(data.Fuel, 'f', precision, 64)
				g.flaps.Text = strconv.FormatFloat(data.Flaps, 'f', precision, 64)
				g.craft.Text = data.Type
				g.w.Invalidate()
				// l.Log("indicator", data)
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
					layout.Spacer{Height: unit.Dp(250)}.Layout,
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(g.th, unit.Sp(30), state.IASKmH)
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						g.ias.Alignment = text.Middle
						return g.ias.Layout(gtx)
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
						lbl := material.Label(g.th, unit.Sp(30), "water temperature")
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						g.waterTemp.Alignment = text.Middle
						return g.waterTemp.Layout(gtx)
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
						lbl := material.Label(g.th, unit.Sp(30), "throttle")
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						g.throttle.Alignment = text.Middle
						return g.throttle.Layout(gtx)
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

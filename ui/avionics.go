package ui

import (
	"context"
	"fmt"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/grafov/kiwi"
	"github.com/wt-tools/wtscope/input/indicators"
	"github.com/wt-tools/wtscope/input/state"
)

const noAircraft = "wait for flight"

type avionics struct {
	w            *app.Window
	th           *material.Theme
	ias          *intBasicDisplay
	altH         *intBasicDisplay
	oilTemp      *floatBasicDisplay
	headTemp     *floatBasicDisplay
	waterTemp    *floatBasicDisplay
	fuel         *floatBasicDisplay
	flaps        *floatBasicDisplay
	throttle     *intBasicDisplay
	craft        material.ButtonStyle
	btnClickArea widget.Clickable
	log          *kiwi.Logger
}

func newAvionics(log *kiwi.Logger) *avionics {
	return &avionics{
		w:   app.NewWindow(app.Title("WT Scope: Avionics")),
		th:  material.NewTheme(gofont.Collection()),
		log: log,
	}
}

const precision = 1 // numbers after comma for floating values
func (g *gui) UpdateAvionics(ctx context.Context, states *state.Service, inds *indicators.Service) {
	l := g.log.New()
	g.av.ias = newIntBasicDisplay(g.av.th, "speed", 350)
	g.av.altH = newIntBasicDisplay(g.av.th, "altitude", 90)
	g.av.oilTemp = newFloatBasicDisplay(g.av.th, "oil temperature", 120)
	g.av.waterTemp = newFloatBasicDisplay(g.av.th, "water temperature", 120)
	g.av.headTemp = newFloatBasicDisplay(g.av.th, "head temperature", 120)
	g.av.fuel = newFloatBasicDisplay(g.av.th, "fuel", 60)
	g.av.flaps = newFloatBasicDisplay(g.av.th, "flaps", 60)
	g.av.throttle = newIntBasicDisplay(g.av.th, "throttle", 60)
	g.av.craft = material.Button(g.av.th, &g.av.btnClickArea, noAircraft)
	go func() {
		for {
			select {
			case data := <-states.Messages:
				g.av.altH.Set(data.GetInt(state.HM))
				g.av.ias.Set(data.GetInt(state.IASKmH))
				g.av.throttle.Set(data.GetInt(state.Throttle1))
				g.av.w.Invalidate()
				l.Log("state", data)
			case data := <-inds.Messages:
				if data.OilTemperature >= 0 {
					g.av.oilTemp.Set(data.OilTemperature)
				}
				if data.HeadTemperature >= 0 {
					g.av.headTemp.Set(data.HeadTemperature)
				}
				if data.WaterTemperature >= 0 {
					g.av.waterTemp.Set(data.WaterTemperature)
				}
				g.av.fuel.Set(data.Fuel)
				g.av.flaps.Set(data.Flaps)
				g.av.craft.Text = data.Type
				g.av.w.Invalidate()
				l.Log("indicator", data)
			}
		}
	}()
}

// TODO split to windows
func (a *avionics) panel() error {
	l := a.log.New()
	var ops op.Ops
	btn1 := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			return a.craft.Layout(gtx)
		},
	)

	btn2 := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(a.th, &a.btnClickArea, "dev preview")

			return btn.Layout(gtx)
		},
	)

	rows := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}

	for e := range a.w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
			visible := !(a.craft.Text == noAircraft) || a.craft.Text == ""
			for _, event := range gtx.Events(a.w) {
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
				layout.Rigid(a.ias.Display(gtx, visible)),
				layout.Rigid(a.altH.Display(gtx, visible)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(0.33, a.oilTemp.Display(gtx, visible)),
						layout.Flexed(0.33, a.waterTemp.Display(gtx, visible)),
						layout.Flexed(0.33, a.headTemp.Display(gtx, visible)),
					)
				}),
				layout.Rigid(a.throttle.Display(gtx, visible)),
				layout.Rigid(a.flaps.Display(gtx, visible)),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return rows.Layout(gtx, btn1, btn2)
					},
				),
			)
			if a.btnClickArea.Clicked() {
				fmt.Println("button was clicked")
			}
			area.Pop()
			e.Frame(gtx.Ops)
		}
	}

	return nil
}

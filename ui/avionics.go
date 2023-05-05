package ui

import (
	"context"
	"fmt"
	"strconv"

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
	g.av.ias = newBasicDisplay(g.av.th, "speed", 330)
	g.av.altH = newBasicDisplay(g.av.th, "altitude", 50)
	g.av.oilTemp = newBasicDisplay(g.av.th, "oil temperature", 90)
	g.av.waterTemp = newBasicDisplay(g.av.th, "water temperature", 90)
	g.av.headTemp = newBasicDisplay(g.av.th, "head temperature", 90)
	g.av.fuel = newBasicDisplay(g.av.th, "fuel", 50)
	g.av.flaps = newBasicDisplay(g.av.th, "flaps", 50)
	g.av.throttle = newBasicDisplay(g.av.th, "throttle", 50)
	g.av.craft = material.Button(g.av.th, &g.av.btnClickArea, noAircraft)
	go func() {
		for {
			select {
			case data := <-states.Messages:
				g.av.altH.V = strconv.Itoa(data.GetInt(state.HM))
				g.av.ias.V = strconv.Itoa(data.GetInt(state.IASKmH))
				g.av.throttle.V = strconv.Itoa(data.GetInt(state.Throttle1))
				g.av.w.Invalidate()
				l.Log("state", data)
			case data := <-inds.Messages:
				if data.OilTemperature >= 0 {
					g.av.oilTemp.V = strconv.FormatFloat(data.OilTemperature, 'f', precision, 64)
				}
				if data.HeadTemperature >= 0 {
					g.av.headTemp.V = strconv.FormatFloat(data.HeadTemperature, 'f', precision, 64)
				}
				if data.WaterTemperature >= 0 {
					g.av.waterTemp.V = strconv.FormatFloat(data.WaterTemperature, 'f', precision, 64)
				}
				g.av.fuel.V = strconv.FormatFloat(data.Fuel, 'f', precision, 64)
				g.av.flaps.V = strconv.FormatFloat(data.Flaps, 'f', precision, 64)
				g.av.craft.Text = data.Type
				g.av.w.Invalidate()
				l.Log("indicator", data)
			}
		}
	}()
}

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
				layout.Rigid(a.oilTemp.Display(gtx, visible)),
				layout.Rigid(a.waterTemp.Display(gtx, visible)),
				layout.Rigid(a.headTemp.Display(gtx, visible)),
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

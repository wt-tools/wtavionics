package ui

import (
	"context"
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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
	compass      *compassDisplay
	ias          *intBasicDisplay
	iasMph       *intBasicDisplay
	iasKnot      *floatBasicDisplay
	altH         *intBasicDisplay
	vsi          *floatBasicDisplay
	oilTemp      *floatBasicDisplay
	headTemp     *floatBasicDisplay
	waterTemp    *floatBasicDisplay
	fuel         *floatBasicDisplay
	flaps        *floatBasicDisplay
	gears        *floatBasicDisplay
	throttle     *intBasicDisplay
	mixture      *intBasicDisplay
	rpmThrottle  *intBasicDisplay
	craft        material.ButtonStyle
	btnClickArea widget.Clickable
	log          *kiwi.Logger
}

func newAvionics(log *kiwi.Logger) *avionics {
	var w app.Window
	w.Option(app.Title("WT Scope: Avionics"))
	return &avionics{
		w:   &w,
		th:  material.NewTheme(),
		log: log,
	}
}

const (
	precision = 1 // numbers after comma for floating values
	knotMul   = 0.539956803
	mphDiv    = 1.609344
)

func (g *gui) UpdateAvionics(ctx context.Context, states *state.Service, inds *indicators.Service) {
	l := g.log.New()
	g.av.compass = newCompassDisplay(g.av.th)
	g.av.ias = newIntBasicDisplay(g.av.th, "speed, KM/h", 320)
	g.av.iasMph = newIntBasicDisplay(g.av.th, "speed, MPH", 120)
	g.av.iasKnot = newFloatBasicDisplay(g.av.th, "speed, Knots", 120)
	g.av.altH = newIntBasicDisplay(g.av.th, "altitude", 120)
	g.av.vsi = newFloatBasicDisplay(g.av.th, "vertical speed", 120)
	g.av.oilTemp = newFloatBasicDisplay(g.av.th, "oil temperature", 150)
	g.av.waterTemp = newFloatBasicDisplay(g.av.th, "water temperature", 150)
	g.av.headTemp = newFloatBasicDisplay(g.av.th, "head temperature", 150)
	g.av.fuel = newFloatBasicDisplay(g.av.th, "fuel", 90)
	g.av.flaps = newFloatBasicDisplay(g.av.th, "flaps", 90)
	g.av.gears = newFloatBasicDisplay(g.av.th, "gears", 90)
	g.av.throttle = newIntBasicDisplay(g.av.th, "throttle", 90)
	g.av.rpmThrottle = newIntBasicDisplay(g.av.th, "RPM Throttle", 90)
	g.av.mixture = newIntBasicDisplay(g.av.th, "mixture", 90)
	g.av.craft = material.Button(g.av.th, &g.av.btnClickArea, noAircraft)
	go func() {
		for {
			select {
			// from state
			case data := <-states.Messages:
				g.av.altH.Set(data.GetInt(state.HM))
				ias := data.GetInt(state.IASKmH)
				g.av.ias.Set(ias)
				g.av.iasKnot.Set(float64(ias) * knotMul)
				g.av.iasMph.Set(int(float64(ias) / mphDiv))
				v := data.GetFloat64(state.VyMS)
				// stop flickering with minuses
				if v > -0.01 && v < 0.01 {
					v = 0
				}
				g.av.vsi.Set(v)
				g.av.throttle.Set(data.GetInt(state.Throttle1))
				g.av.w.Invalidate()
				l.Log("state", data)
			// from indicators
			case data := <-inds.Messages:
				g.av.compass.Set(int(data.Compass))
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
				g.av.mixture.Set(int(data.Mixture1)) // TODO fix it for multiple engines
				g.av.flaps.Set(data.Flaps)
				g.av.gears.Set(data.Gears)
				g.av.craft.Text = data.Type
				g.av.w.Invalidate()
				l.Log("indicator", data)
			}
		}
	}()
}

// XXX doesn't work yet :\
func (a *avionics) exitOnEsc(gtx layout.Context, tag bool) {
	event.Op(gtx.Ops, tag)
	// New event reading
	for {
		event, ok := gtx.Event(
			key.FocusFilter{
				Target: tag,
			},
			key.Filter{
				Focus: tag,
				Name:  key.NameEscape,
			},
			key.Filter{
				Focus: tag,
				Name:  key.NameEnter,
			},
		)
		if !ok {
			break
		}
		ev, ok := event.(key.Event)
		if !ok {
			continue
		}
		if ev.Name == key.NameEscape {
			os.Exit(0)
		}
	}
}

// TODO split to windows
func (a *avionics) panel() error {
	// l := a.log.New()
	var ops op.Ops
	btn1 := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			return a.craft.Layout(gtx)
		},
	)

	btn2 := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			return material.Button(a.th, &a.btnClickArea, "dev preview").Layout(gtx)
		},
	)

	rows := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}
	var exitTag bool
	for {
		switch e := a.w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.FillShape(gtx.Ops, bgColor, clip.Rect{Max: gtx.Constraints.Max}.Op())
			area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
			visible := !(a.craft.Text == noAircraft) || a.craft.Text == ""
			a.exitOnEsc(gtx, exitTag)
			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				layout.Rigid(
					layout.Spacer{Height: unit.Dp(25)}.Layout,
				),
				layout.Rigid(a.compass.Display(gtx, visible)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Spacing: layout.SpaceEnd,
					}.Layout(gtx,
						layout.Flexed(0.25, a.iasMph.Display(gtx, visible)),
						layout.Flexed(0.5, a.ias.Display(gtx, visible)),
						layout.Flexed(0.25, a.iasKnot.Display(gtx, visible)),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(0.5, a.altH.Display(gtx, visible)),
						layout.Flexed(0.5, a.vsi.Display(gtx, visible)),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(0.333, a.oilTemp.Display(gtx, visible)),
						layout.Flexed(0.333, a.waterTemp.Display(gtx, visible)),
						layout.Flexed(0.333, a.headTemp.Display(gtx, visible)),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(0.333, a.throttle.Display(gtx, visible)),
						layout.Flexed(0.333, a.mixture.Display(gtx, visible)),
						layout.Flexed(0.333, a.fuel.Display(gtx, visible)),
					)
				}),
				layout.Rigid(a.flaps.Display(gtx, visible)),
				layout.Rigid(a.gears.Display(gtx, visible)),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return rows.Layout(gtx, btn1, btn2)
					},
				),
			)
			if a.btnClickArea.Clicked(gtx) {
				fmt.Println("button was clicked")
			}
			area.Pop()
			e.Frame(gtx.Ops)
		}
	}
}

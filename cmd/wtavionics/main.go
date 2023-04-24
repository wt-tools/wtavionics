package main

import (
	"context"
	"fmt"
	"image/color"
	"net/http"
	"os"

	"github.com/wt-tools/wtavionics/config"
	"github.com/wt-tools/wtscope/input/indicators"
	"github.com/wt-tools/wtscope/input/state"
	"github.com/wt-tools/wtscope/net/poll"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/grafov/kiwi"
)

func main() {
	ctx := context.Background()
	kiwi.SinkTo(os.Stdout, kiwi.AsLogfmt()).Start()
	l := kiwi.New()
	conf := config.New()
	l.Log("status", "prepare avionics for start", "config", "xxx")
	errch := make(chan error, 8) // XXX разделить по компонентам
	defaultPolling := poll.New(http.DefaultClient, errch, 2, 2)
	stateSvc := state.New(conf, defaultPolling, errch)
	indSvc := indicators.New(conf, defaultPolling, errch)
	go defaultPolling.Do()
	go stateSvc.Grab(ctx)
	go indSvc.Grab(ctx)
	go func() {
		w := app.NewWindow(app.Title("WT Scope: Avionics"))
		err := run(w, l, stateSvc, indSvc)
		if err != nil {
			l.Log("fatal", "can't run window", "error", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window, l *kiwi.Logger, states *state.Service, inds *indicators.Service) error {
	th := material.NewTheme(gofont.Collection())
	speed := material.H1(th, "0")
	go func() {
		for {
			select {
			case data := <-states.Messages:
				l.Log("state", data)
			case data := <-inds.Messages:
				speed.Text = fmt.Sprintf("%f\n%f\n%f\n%f", data.Speed, data.AltitudeHour, data.OilTemperature, data.HeadTemperature)
				w.Invalidate()
				l.Log("indicator", data)
			}
		}
	}()
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
			key.InputOp{Tag: w, Keys: key.NameEscape}.Add(gtx.Ops)
			key.FocusOp{Tag: w}.Add(gtx.Ops)
			for _, event := range gtx.Events(w) {
				switch event := event.(type) {
				case key.Event:
					l.Log("event", "key pressed")
					if event.Name == key.NameEscape {
						return nil
					}
				}
			}
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			speed.Color = maroon
			speed.Alignment = text.Middle
			speed.Layout(gtx)
			area.Pop()
			e.Frame(gtx.Ops)
		}
	}
}

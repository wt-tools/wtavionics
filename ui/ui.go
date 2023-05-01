package ui

import (
	"context"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/grafov/kiwi"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type gui struct {
	w   *app.Window
	th  *material.Theme
	log *kiwi.Logger
	avionicsDisplays
}

func Init(_ context.Context, log *kiwi.Logger) *gui {
	return &gui{
		w:   app.NewWindow(app.Title("WT Scope: Avionics")),
		th:  material.NewTheme(gofont.Collection()),
		log: log,
	}
}

func (g *gui) Run(_ context.Context) {
	l := g.log.New()
	go func() {
		err := g.avionicsPanel()
		if err != nil {
			l.Log("fatal", "can't run window", "error", err)
		}
		l.Log("exit", "exit by escape")
		os.Exit(0)
	}()
	app.Main()
}

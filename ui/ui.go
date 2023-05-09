package ui

import (
	"context"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"github.com/grafov/kiwi"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type gui struct {
	log *kiwi.Logger
	av  *avionics
}

func Init(_ context.Context, log *kiwi.Logger) *gui {
	return &gui{
		log: log,
		av:  newAvionics(log),
	}
}

func (g *gui) Run(_ context.Context) {
	l := g.log.New()
	go func() {
		err := g.av.panel()
		if err != nil {
			l.Log("fatal", "can't run avionics window", "error", err)
			os.Exit(0)
		}
	}()
	app.Main()
}

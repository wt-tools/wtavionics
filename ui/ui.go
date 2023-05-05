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
	th  *material.Theme
	log *kiwi.Logger
	av  *avionics
	bl  *battleLog
}

func Init(_ context.Context, log *kiwi.Logger) *gui {
	th := material.NewTheme(gofont.Collection())
	return &gui{
		th:  th,
		log: log,
		av:  newAvionics(th, log),
		bl:  newBattleLog(th, log),
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
	go func() {
		err := g.bl.panel()
		if err != nil {
			l.Log("fatal", "can't run battle log window", "error", err)
			os.Exit(0)
		}
	}()
	app.Main()
}

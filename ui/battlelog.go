package ui

import (
	"context"
	"fmt"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/grafov/kiwi"
	"github.com/wt-tools/wtscope/action"
	"github.com/wt-tools/wtscope/input/hudmsg"
)

var headings []string

type battleLog struct {
	w          *app.Window
	th         *material.Theme
	log        *kiwi.Logger
	grid       component.GridState
	list       widget.List
	rows       []action.GeneralAction
	latestTime time.Duration
}

func newBattleLog(log *kiwi.Logger) *battleLog {
	return &battleLog{
		w:   app.NewWindow(app.Title("WT Scope: Battle Log")),
		th:  material.NewTheme(gofont.Collection()),
		log: log,
	}
}

func (g *gui) UpdateBattleLog(ctx context.Context, gamelog *hudmsg.Service) {
	l := g.log.New()
	go func() {
		for {
			select {
			case data := <-gamelog.Messages:
				if len(g.bl.rows) > 0 && g.bl.latestTime > data.At {
					g.bl.rows = nil
				}
				g.bl.latestTime = data.At
				g.bl.rows = append(g.bl.rows, data)
				g.bl.w.Invalidate()
				l.Log("battle log", data)
			}
		}
	}()
}

func (b *battleLog) panel() error {
	var ops op.Ops
	b.list.Axis = layout.Vertical
	b.list.ScrollToEnd = true
	for {
		e := <-b.w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			if len(b.rows) == 0 {
				continue
			}
			gtx := layout.NewContext(&ops, e)
			b.listLayout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (b *battleLog) listLayout(gtx C) D {
	return material.List(b.th, &b.list).Layout(gtx, len(b.rows), func(gtx layout.Context, i int) layout.Dimensions {
		var text string
		switch {
		case len(b.rows) == 0:
			text = "no battle log yet"
		case i > len(b.rows):
			text = fmtAction(b.rows[len(b.rows)-1])
		default:
			text = fmtAction(b.rows[i])
		}
		return material.Label(b.th, unit.Sp(26), text).Layout(gtx)
	})
}

func fmtAction(a action.GeneralAction) string {
	return fmt.Sprintf("%16s  %s", a.At, a.Origin)
}

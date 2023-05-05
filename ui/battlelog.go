package ui

import (
	"context"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/grafov/kiwi"
	"github.com/wt-tools/wtscope/action"
	"github.com/wt-tools/wtscope/input/hudmsg"
)

var headings = []string{"At", "Message", "Attacker", "Vehicle", "Target player", "Target vehicle"}

type battleLog struct {
	w    *app.Window
	th   *material.Theme
	log  *kiwi.Logger
	grid component.GridState
	rows []action.GeneralAction
}

func newBattleLog(th *material.Theme, log *kiwi.Logger) *battleLog {
	return &battleLog{
		w:   app.NewWindow(app.Title("WT Scope: Battle Log")),
		th:  th,
		log: log,
	}
}

func (g *gui) UpdateBattleLog(ctx context.Context, gamelog *hudmsg.Service) {
	l := g.log.New()
	go func() {
		for {
			select {
			case data := <-gamelog.Messages:
				// msg := data.Origin
				g.bl.rows = append(g.bl.rows, data)
				g.bl.w.Invalidate()
				l.Log("battle log", data)
			}
		}
	}()
}

func (b *battleLog) panel() error {
	var ops op.Ops
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
			b.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (b *battleLog) layout(gtx C) D {
	// Configure width based on available space and a minimum size.
	minSize := gtx.Dp(unit.Dp(800))
	border := widget.Border{
		Color: color.NRGBA{A: 255},
		Width: unit.Dp(1),
	}

	inset := layout.UniformInset(unit.Dp(2))

	// Configure a label styled to be a heading.
	headingLabel := material.Body1(b.th, "")
	headingLabel.Font.Weight = font.Bold
	headingLabel.Alignment = text.Middle
	headingLabel.MaxLines = 1

	// Configure a label styled to be a data element.
	dataLabel := material.Body1(b.th, "")
	dataLabel.Font.Variant = "Mono"
	dataLabel.MaxLines = 1
	dataLabel.Alignment = text.End

	// Measure the height of a heading row.
	orig := gtx.Constraints
	gtx.Constraints.Min = image.Point{}
	macro := op.Record(gtx.Ops)
	dims := inset.Layout(gtx, headingLabel.Layout)
	_ = macro.Stop()
	gtx.Constraints = orig

	return component.Table(b.th, &b.grid).Layout(gtx, len(b.rows), len(headings),
		func(axis layout.Axis, index, constraint int) int {
			widthUnit := max(int(float32(constraint)/2), minSize)
			switch axis {
			case layout.Horizontal:
				switch index {
				case 0:
					return int(widthUnit / 8)
				case 1:
					return int(widthUnit)
				case 2, 3, 4, 5, 6:
					return int(widthUnit / 5)
				default:
					return 0
				}
			default:
				return dims.Size.Y
			}
		},
		func(gtx C, col int) D {
			return border.Layout(gtx, func(gtx C) D {
				return inset.Layout(gtx, func(gtx C) D {
					headingLabel.Text = headings[col]
					return headingLabel.Layout(gtx)
				})
			})
		},
		func(gtx C, row, col int) D {
			return inset.Layout(gtx, func(gtx C) D {
				switch col {
				case 0:
					dataLabel.Text = b.rows[row].At.Format("03:04:05")
				case 1:
					dataLabel.Text = b.rows[row].Origin
				case 2:
					dataLabel.Text = b.rows[row].Damage.Player.Name
				case 3:
					dataLabel.Text = b.rows[row].Damage.Vehicle.Name
				case 4:
					dataLabel.Text = b.rows[row].Damage.TargetPlayer.Name
				case 5:
					dataLabel.Text = b.rows[row].Damage.TargetVehicle.Name
				default:
					dataLabel.Text = "unknown value"

				}
				return dataLabel.Layout(gtx)
			})
		},
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

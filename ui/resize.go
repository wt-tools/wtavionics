package ui

import (
	"image"

	"gioui.org/unit"
)

func dp(size image.Point) unit.Sp {
	return unit.Sp(5 + 0.01875*float32(size.X-320)) // linear approximation
}

func height(size image.Point, h int) unit.Sp {
	var w unit.Sp
	switch {
	case size.X > 1600:
		w = 1
	case size.X > 1200:
		w = 0.7
	case size.X > 800:
		w = 0.5
	case size.X > 640:
		w = 0.3
	case size.X > 320:
		w = 0.2
	}
	return w * unit.Sp(h)
}

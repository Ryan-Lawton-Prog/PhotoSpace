package utils

import (
	"image"

	"gioui.org/layout"
)

func CombineLayout(a, b layout.Dimensions) layout.Dimensions {
	return layout.Dimensions{
		Size: image.Point{
			X: max(a.Size.X, b.Size.X),
			Y: max(a.Size.Y, b.Size.Y),
		},
		Baseline: max(a.Baseline, b.Baseline),
	}
}

func AddLayoutHeight(a, b layout.Dimensions) layout.Dimensions {
	return layout.Dimensions{
		Size: image.Point{
			X: max(a.Size.X, b.Size.X),
			Y: a.Size.Y + b.Size.Y,
		},
		Baseline: max(a.Baseline, b.Baseline),
	}
}

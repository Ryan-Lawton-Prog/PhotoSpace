package shapes

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace-ui/models"
	"ryanlawton.art/photospace-ui/widgets/colors"
)

type Params struct {
	Color     color.NRGBA
	Alignment text.Alignment
	Text      string
	Theme     *material.Theme
	Size      Size
	Shadow    bool
}

type draw func(*models.C, Params) func()

type Offset [2]int
type Size struct {
	Width  int
	Height int
}

// color.NRGBA{R: 0x80, A: 0xFF}
func DrawSquare(gtx *models.C, params Params) func() {
	rect := clip.RRect{image.Rectangle{Max: image.Pt(params.Size.Width, params.Size.Height)}, 20, 20, 20, 20}.Push(gtx.Ops)
	paint.ColorOp{Color: params.Color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return rect.Pop
}

func OffsetDraw(gtx *models.C, offset Offset, draw draw, params Params) {
	defer op.Offset(image.Pt(offset[0], offset[1])).Push(gtx.Ops).Pop()
	draw(gtx, params)
}

func DrawText(gtx *models.C, params Params) layout.Dimensions {
	var shadow material.LabelStyle

	var shadowLayout layout.Dimensions
	var layout layout.Dimensions

	if params.Shadow {
		offset := image.Pt(0, 20)
		offsetPop := op.Offset(offset).Push(gtx.Ops)
		shadow = material.H1(params.Theme, params.Text)
		shadow.Color = colors.MainTheme.Gray0
		shadow.Font.Weight = font.Bold
		shadow.Alignment = params.Alignment
		shadowLayout.Baseline = shadow.Layout(*gtx).Baseline
		shadowLayout.Size = offset

		layout = combineLayout(layout, shadowLayout)

		offsetPop.Pop()
	}

	// Define an large label with an appropriate text:
	title := material.H1(params.Theme, params.Text)
	title.Color = params.Color
	title.Font.Weight = font.Bold

	// Change the position of the label.
	title.Alignment = params.Alignment

	layout = combineLayout(title.Layout(*gtx), layout)

	// Draw the label to the graphics context.
	return layout
}

func combineLayout(a, b layout.Dimensions) layout.Dimensions {
	return layout.Dimensions{
		Size: image.Point{
			X: max(a.Size.X, b.Size.X),
			Y: max(a.Size.Y, b.Size.Y),
		},
		Baseline: max(a.Baseline, b.Baseline),
	}
}

func DrawTextBox(gtx *models.C, tb *widget.Editor, params Params) layout.Dimensions {
	// Wrap the editor in material design
	textBox := material.Editor(params.Theme, tb, params.Text)
	// Define characteristics of the input box
	tb.SingleLine = true
	tb.Alignment = text.Middle

	l := textBox.Layout(*gtx)

	defer DrawSquare(gtx, Params{
		Size: Size{
			Width:  l.Size.X,
			Height: l.Size.Y,
		},
		Color: colors.MainTheme.Warning,
	})()

	// Define insets ...
	// margins := layout.Inset{
	// 	Top:    unit.Dp(0),
	// 	Right:  unit.Dp(170),
	// 	Bottom: unit.Dp(40),
	// 	Left:   unit.Dp(170),
	// }
	// ... and borders ...
	border := widget.Border{
		Color:        colors.Red,
		CornerRadius: unit.Dp(10),
		Width:        unit.Dp(1),
	}
	// ... before laying it out, one inside the other
	return border.Layout(*gtx, textBox.Layout)
}

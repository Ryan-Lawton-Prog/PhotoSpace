package shapes

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace-ui/models"
)

type Params struct {
	Color     color.NRGBA
	Alignment text.Alignment
	Text      string
	Theme     *material.Theme
	Size      Size
}

type draw func(*models.C, Params)

type Offset [2]int
type Size struct {
	Width  int
	Height int
}

// color.NRGBA{R: 0x80, A: 0xFF}
func DrawSquare(gtx *models.C, params Params) {
	rect := clip.Rect{Max: image.Pt(params.Size.Width, params.Size.Height)}.Push(gtx.Ops)
	paint.ColorOp{Color: params.Color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	rect.Pop()
}

func OffsetDraw(gtx *models.C, offset Offset, draw draw, params Params) {
	defer op.Offset(image.Pt(offset[0], offset[1])).Push(gtx.Ops).Pop()
	draw(gtx, params)
}

func DrawText(gtx *models.C, params Params) layout.Dimensions {
	// Define an large label with an appropriate text:
	title := material.H1(params.Theme, params.Text)
	title.Color = params.Color

	// Change the position of the label.
	title.Alignment = params.Alignment

	// Draw the label to the graphics context.
	return title.Layout(*gtx)
}

func DrawTextBox(gtx *models.C, tb *widget.Editor, params Params) layout.Dimensions {
	// Wrap the editor in material design
	textBox := material.Editor(params.Theme, tb, params.Text)
	// Define characteristics of the input box
	tb.SingleLine = true
	tb.Alignment = text.Middle

	// Define insets ...
	margins := layout.Inset{
		Top:    unit.Dp(0),
		Right:  unit.Dp(170),
		Bottom: unit.Dp(40),
		Left:   unit.Dp(170),
	}
	// ... and borders ...
	border := widget.Border{
		Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
		CornerRadius: unit.Dp(3),
		Width:        unit.Dp(2),
	}
	// ... before laying it out, one inside the other
	return margins.Layout(*gtx,
		func(gtx models.C) models.D {
			return border.Layout(gtx, textBox.Layout)
		},
	)
}

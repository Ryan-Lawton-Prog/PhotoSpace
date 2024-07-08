package shapes

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace/internal/app/models"
	"ryanlawton.art/photospace/internal/pkg/themes"
)

type SquareParams struct {
	Color     color.NRGBA
	Alignment text.Alignment
	Text      string
	Theme     *material.Theme
	Size      Size
	Shadow    bool
	Gradient  bool
}

func DrawSquare(gtx *models.C, params SquareParams) func() {
	rect := clip.RRect{
		Rect: image.Rectangle{
			Max: image.Pt(params.Size.Width, params.Size.Height),
		},
		SE: 1,
		SW: 1,
		NW: 1,
		NE: 1,
	}.Push(gtx.Ops)
	if params.Gradient {
		paint.LinearGradientOp{
			Stop1:  f32.Pt(float32(gtx.Constraints.Max.X)/2, float32(params.Size.Height)),
			Color1: themes.White,
			Stop2:  f32.Pt(float32(gtx.Constraints.Max.X)/2, 0),
			Color2: themes.MainTheme.Gray1,
		}.Add(gtx.Ops)
	} else {
		paint.ColorOp{Color: params.Color}.Add(gtx.Ops)
	}

	paint.PaintOp{}.Add(gtx.Ops)

	return rect.Pop
}

func DrawSnakeSquare(gtx *models.C, params Params) func() {
	rect := clip.RRect{image.Rectangle{Max: image.Pt(params.Size.Width, params.Size.Height)}, 1, 1, 1, 1}.Push(gtx.Ops)
	if params.Gradient {
		paint.LinearGradientOp{
			Stop1:  f32.Pt(float32(gtx.Constraints.Max.X)/2, float32(params.Size.Height)),
			Color1: themes.White,
			Stop2:  f32.Pt(float32(gtx.Constraints.Max.X)/2, 0),
			Color2: themes.MainTheme.Gray1,
		}.Add(gtx.Ops)
	} else {
		paint.ColorOp{Color: params.Color}.Add(gtx.Ops)
	}

	paint.PaintOp{}.Add(gtx.Ops)

	return rect.Pop
}

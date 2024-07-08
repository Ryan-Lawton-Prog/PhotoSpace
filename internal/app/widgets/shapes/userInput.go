package shapes

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace/internal/app/utils"
	"ryanlawton.art/photospace/internal/pkg/themes"
)

func DrawFormTextBox(tb *widget.Editor, params Params) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		var textBoxD layout.Dimensions
		// Wrap the editor in material design
		textBox := material.Editor(params.Theme, tb, params.Text)
		// Define characteristics of the input box
		tb.SingleLine = true
		tb.Alignment = text.Middle
		tb.LineHeight = unit.Sp(30)
		textBox.LineHeight = unit.Sp(50)
		// todo: text size constants
		textBox.TextSize = 24
		tb.LineHeightScale = 5

		box := textBox.Layout(gtx)

		DrawSquare(&gtx, SquareParams{
			Size: Size{
				Width:  box.Size.X,
				Height: box.Size.Y,
			},
			Color: themes.MainTheme.Gray0,
		})()

		if params.Shadow {
			offset := image.Pt(0, box.Size.Y-(2*int(gtx.Metric.PxPerDp)))
			offsetPop := op.Offset(offset).Push(gtx.Ops)
			DrawSquare(&gtx, SquareParams{
				Size: Size{
					Width:  box.Size.X,
					Height: 10 * int(gtx.Metric.PxPerDp),
				},
				Color:    themes.MainTheme.Gray1,
				Gradient: true,
			})()
			b2 := widget.Border{
				Color:        themes.White,
				CornerRadius: unit.Dp(2),
				Width:        unit.Dp(2),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{
					Size: image.Point{
						X: box.Size.X,
						Y: 10 * int(gtx.Metric.PxPerDp),
					},
				}
			})
			offsetPop.Pop()
			textBoxD.Size = textBoxD.Size.Add(image.Pt(0, b2.Size.Y))
		}

		border := widget.Border{
			Color:        themes.MainTheme.Gray10,
			CornerRadius: unit.Dp(1),
			Width:        unit.Dp(2),
		}
		// ... before laying it out, one inside the other
		textBoxD = utils.AddLayoutHeight(border.Layout(gtx, textBox.Layout), textBoxD)

		return textBoxD
	}
}

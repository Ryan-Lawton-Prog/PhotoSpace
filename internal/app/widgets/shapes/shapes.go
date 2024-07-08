package shapes

import (
	"image"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace/internal/app/models"
	"ryanlawton.art/photospace/internal/app/utils"
	"ryanlawton.art/photospace/internal/pkg/themes"
)

func DrawTextBox(gtx *models.C, tb *widget.Editor, params Params) layout.Dimensions {
	var layout layout.Dimensions
	// Wrap the editor in material design
	textBox := material.Editor(params.Theme, tb, params.Text)
	// Define characteristics of the input box
	tb.SingleLine = true
	tb.Alignment = text.Middle

	l := textBox.Layout(*gtx)

	defer DrawSquare(gtx, SquareParams{
		Size: Size{
			Width:  l.Size.X,
			Height: l.Size.Y,
		},
		Color:  themes.MainTheme.Gray2,
		Shadow: params.Shadow,
	})()

	border := widget.Border{
		Color:        themes.MainTheme.Gray10,
		CornerRadius: unit.Dp(1),
		Width:        unit.Dp(1),
	}
	// ... before laying it out, one inside the other
	layout = utils.CombineLayout(border.Layout(*gtx, textBox.Layout), layout)
	if params.Shadow {
		layout.Size = layout.Size.Add(image.Pt(0, 20))
	}

	return layout
}

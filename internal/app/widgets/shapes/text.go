package shapes

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace/internal/app/models"
	"ryanlawton.art/photospace/internal/app/utils"
	"ryanlawton.art/photospace/internal/pkg/themes"
)

type Params struct {
	Color     color.NRGBA
	Alignment text.Alignment
	Text      string
	Theme     *material.Theme
	Size      Size
	Shadow    bool
	Gradient  bool
}

type TextParams struct {
	Color     color.NRGBA
	Alignment text.Alignment
	Text      string
	Theme     *material.Theme
	Size      TextSize
	Shadow    bool
}

type draw func(*models.C, Params) func()

type Offset [2]int

func OffsetDraw(gtx *models.C, offset Offset, draw draw, params Params) {
	defer op.Offset(image.Pt(offset[0], offset[1])).Push(gtx.Ops).Pop()
	draw(gtx, params)
}

type TextSize string

const (
	H1 TextSize = "H1"
	H2 TextSize = "H2"
	H3 TextSize = "H3"
	H4 TextSize = "H4"
)

var sizeMap = map[TextSize]func(th *material.Theme, txt string) material.LabelStyle{
	H1: material.H1,
	H2: material.H2,
	H3: material.H3,
	H4: material.H4,
}

// todo: move to different file and fix shadow outbounds
func DrawText(params TextParams) layout.Widget {
	return func(gtx models.C) models.D {
		// set default value for text size
		if params.Size == "" {
			params.Size = H1
		}
		var shadow material.LabelStyle

		var shadowLayout layout.Dimensions
		var layout layout.Dimensions

		if params.Shadow {
			offset := image.Pt(0, 5*int(gtx.Metric.PxPerDp))
			offsetPop := op.Offset(offset).Push(gtx.Ops)
			shadow = sizeMap[params.Size](params.Theme, params.Text)
			shadow.Color = themes.MainTheme.Gray0
			shadow.Font.Weight = font.Bold
			shadow.Alignment = params.Alignment
			shadowLayout.Baseline = shadow.Layout(gtx).Baseline
			shadowLayout.Size = offset

			layout = utils.CombineLayout(layout, shadowLayout)

			offsetPop.Pop()
		}

		// Define an large label with an appropriate text:
		title := sizeMap[params.Size](params.Theme, params.Text)
		title.Color = params.Color
		title.Font.Weight = font.Bold

		// Change the position of the label.
		title.Alignment = params.Alignment

		layout = utils.CombineLayout(title.Layout(gtx), layout)

		// Draw the label to the graphics context.
		return layout
	}
}

package loginPage

import (
	"bytes"
	"image"
	"image/jpeg"
	"math"
	"strconv"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace/internal/app/models"
	errorWidget "ryanlawton.art/photospace/internal/app/widgets/error"
	"ryanlawton.art/photospace/internal/app/widgets/shapes"
	"ryanlawton.art/photospace/internal/pkg/logic"
	"ryanlawton.art/photospace/internal/pkg/themes"
)

const (
	minPageWidth  = 400
	sideBarWidth  = 200
	maxImageWidth = 200
	toolBarHeight = 25
)

type Widgets struct {
}

type Image []byte

var list = layout.List{
	Axis: layout.Vertical,
}

type Home struct {
	widgets            Widgets
	pageQueue          chan models.PageId
	errorMessages      errorWidget.ErrorMessages
	errorMessagesQueue chan error
	images             []image.Image
	position           layout.Position
}

func NewPage(pageQueue *chan models.PageId, window *app.Window) *Home {
	widgets := Widgets{}
	ids, _ := logic.GetPhotoIDs()
	images := []image.Image{}
	for _, id := range ids {
		photoB, _ := logic.GetPhotoThumbnail(id)
		file, err := jpeg.Decode(bytes.NewReader(photoB))
		if err != nil {
			continue
		}
		images = append(images, file)
	}
	return &Home{
		widgets:            widgets,
		pageQueue:          *pageQueue,
		errorMessages:      errorWidget.NewErrorMessageWidget(window.Invalidate),
		errorMessagesQueue: make(chan error),
		images:             images,
		position:           layout.Position{},
	}
}

func (page *Home) handleInput(gtx *models.C) {

}

func (page *Home) StartRoutines(window *app.Window) {
	page.errorMessages.SetRefresh(window.Invalidate)
	page.errorMessages.StartRoutine()

	go func() {
		for err := range page.errorMessagesQueue {
			page.errorMessages.Add(err.Error(), time.Now().Add(time.Second*10))
		}
	}()
}

// UploadPhoto uploads a photo to the database
func (page *Home) Layout(gtx *models.C, th *material.Theme) {
	// Handle user input before rendering
	page.handleInput(gtx)

	//left := int(max((float32(gtx.Constraints.Max.X)/gtx.Metric.PxPerDp-minPageWidth)/2, 0))

	// Calculate the margin borders
	inset := layout.Inset{}

	// Layout screen
	inset.Layout(*gtx, func(gtx models.C) models.D {
		flex := layout.Flex{
			// Vertical alignment, from top to bottom
			Axis: layout.Vertical,
			// Empty space is left at the start, i.e. at the top
			Spacing: layout.SpaceEnd,
		}.Layout(gtx,
			// TITLE
			// layout.Rigid(
			// shapes.DrawText(shapes.TextParams{
			// 	Theme:     th,
			// 	Text:      string(page.images[0]),
			// 	Color:     themes.MainTheme.Gray10,
			// 	Alignment: text.Middle,
			// 	Shadow:    true,
			// 	Size:      shapes.H3,
			// }),
			// ),
			page.toolbar(&gtx, th),
			layout.Rigid(
				func(gtx models.C) layout.Dimensions {
					return layout.Flex{
						Axis:    layout.Horizontal,
						Spacing: layout.SpaceBetween,
					}.Layout(
						gtx,
						page.sidebar(&gtx, th),
						page.imageDisplay(),
					)
				},
			),
			page.footer(&gtx, th),
		)

		return flex
	})

}

func (page *Home) toolbar(gtx *models.C, th *material.Theme) layout.FlexChild {
	return layout.Rigid(
		func(gtx models.C) models.D {
			shapes.DrawSquare(&gtx, shapes.SquareParams{
				Color: themes.Red,
				Size: shapes.Size{
					Width:  gtx.Constraints.Max.X,
					Height: int(gtx.Metric.PxPerDp) * toolBarHeight,
				},
			})()

			return layout.Dimensions{
				Size: image.Point{
					Y: int(gtx.Metric.PxPerDp) * toolBarHeight,
				},
			}
		},
	)
}

func (page *Home) sidebar(gtx *models.C, th *material.Theme) layout.FlexChild {
	return layout.Rigid(
		func(gtx models.C) models.D {
			shapes.DrawSquare(&gtx, shapes.SquareParams{
				Color: themes.MainTheme.Error,
				Size: shapes.Size{
					Width:  int(gtx.Metric.PxPerDp) * sideBarWidth,
					Height: gtx.Constraints.Max.Y,
				},
			})()

			return layout.Dimensions{
				Size: image.Point{
					X: int(gtx.Metric.PxPerDp) * sideBarWidth,
					Y: gtx.Constraints.Max.Y,
				},
			}
		},
	)
}

func (page *Home) footer(gtx *models.C, th *material.Theme) layout.FlexChild {
	return layout.Rigid(
		// The height of the spacer is 25 Device independent pixels
		// layout.Spacer{Height: unit.Dp(25)}.Layout,
		shapes.DrawText(shapes.TextParams{
			Theme:     th,
			Text:      strconv.Itoa(len(page.images)),
			Color:     themes.MainTheme.Gray10,
			Alignment: text.Middle,
			Shadow:    true,
			Size:      shapes.H3,
		}),
	)
}

func (page *Home) imageDisplay() layout.FlexChild {
	return layout.Flexed(1,
		func(gtx models.C) models.D {
			adjustedMaxImageWidth := int(maxImageWidth * gtx.Metric.PxPerDp)
			remainingSpace := (gtx.Constraints.Max.X)
			maxWidth := remainingSpace / adjustedMaxImageWidth
			totalImages := len(page.images)
			rows := int(math.Ceil(float64(totalImages) / float64(maxWidth)))

			marginsSize := (gtx.Constraints.Max.X - (adjustedMaxImageWidth * maxWidth)) / 2

			list.Layout(gtx, rows, func(gtx layout.Context, i int) layout.Dimensions {
				for j, im := range page.images[maxWidth*i : min(maxWidth*i+maxWidth, totalImages)] {
					x := im.Bounds().Max
					if im == nil {
						continue
					}

					g := op.Offset(
						image.Pt(
							marginsSize+(j*adjustedMaxImageWidth)%(maxWidth*adjustedMaxImageWidth),
							0),
					).Push(gtx.Ops)
					shapes.DrawSquare(&gtx, shapes.SquareParams{
						Size: shapes.Size{
							Width:  adjustedMaxImageWidth,
							Height: adjustedMaxImageWidth,
						},
						Color: themes.MainTheme.Gray0,
					})()

					widget.Border{
						Color:        themes.MainTheme.Gray10,
						CornerRadius: unit.Dp(1),
						Width:        unit.Dp(2),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{
							Size: image.Point{
								X: adjustedMaxImageWidth,
								Y: adjustedMaxImageWidth,
							},
						}
					})

					shapes.DrawImage(gtx.Ops, im, x, image.Point{X: adjustedMaxImageWidth, Y: adjustedMaxImageWidth})
					g.Pop()
				}

				return layout.Dimensions{Size: image.Point{
					X: adjustedMaxImageWidth * maxWidth,
					Y: adjustedMaxImageWidth,
				}}
			})

			return layout.Dimensions{Size: image.Point{
				X: gtx.Constraints.Max.X,
				Y: gtx.Constraints.Max.Y,
			}}
		},
	)
}

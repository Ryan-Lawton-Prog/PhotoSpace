package loginPage

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
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
	maxImageWidth = 300
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
		photoB, _ := logic.GetPhoto(id)
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
			// 	shapes.DrawText(shapes.TextParams{
			// 		Theme:     th,
			// 		Text:      string(page.images[0]),
			// 		Color:     themes.MainTheme.Gray10,
			// 		Alignment: text.Middle,
			// 		Shadow:    true,
			// 		Size:      shapes.H3,
			// 	}),
			// ),
			page.Toolbar(&gtx, th),
			layout.Rigid(
				func(gtx models.C) layout.Dimensions {
					return layout.Flex{
						Axis:    layout.Horizontal,
						Spacing: layout.SpaceEnd,
					}.Layout(
						gtx,
						page.Sidebar(&gtx, th),
						page.ImageDisplay(&gtx, th),
					)
				},
			),
			// ... then one to hold an empty spacer
			layout.Rigid(
				// The height of the spacer is 25 Device independent pixels
				layout.Spacer{Height: unit.Dp(25)}.Layout,
			),
		)

		return flex
	})

}

func (page *Home) Toolbar(gtx *models.C, th *material.Theme) layout.FlexChild {

	return layout.Rigid(
		// The height of the spacer is 25 Device independent pixels
		layout.Spacer{Height: unit.Dp(25)}.Layout,
	)

}

func (page *Home) Sidebar(gtx *models.C, th *material.Theme) layout.FlexChild {
	return layout.Rigid(
		layout.Spacer{Width: unit.Dp(sideBarWidth)}.Layout,
	)
}

func (page *Home) ImageDisplay(gtx *models.C, th *material.Theme) layout.FlexChild {

	//remainingSpace := (gtx.Constraints.Max.X / int(gtx.Metric.PxPerDp)) - sideBarWidth
	//maxWidth := remainingSpace/maxImageWidth + 1

	return layout.Flexed(100,
		func(gtx models.C) models.D {
			dim := list.Layout(gtx, len(page.images), func(gtx layout.Context, i int) layout.Dimensions {
				im := page.images[i]
				x := im.Bounds().Max
				//fmt.Println(point, i, x)
				if im == nil {
					return layout.Dimensions{}
				}

				fmt.Println(i)

				// g := op.Offset(
				// 	image.Pt(
				// 		((i * maxImageWidth) % (maxWidth * maxImageWidth)),
				// 		((i)/maxWidth)*maxImageWidth),
				// ).Push(gtx.Ops)
				shapes.DrawSquare(&gtx, shapes.SquareParams{
					Size: shapes.Size{
						Width:  300,
						Height: 300,
					},
					Color: themes.MainTheme.Gray0,
				})()

				size := widget.Border{
					Color:        themes.MainTheme.Gray10,
					CornerRadius: unit.Dp(1),
					Width:        unit.Dp(2),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{
						Size: image.Point{
							X: maxImageWidth,
							Y: maxImageWidth,
						},
					}
				}).Size

				shapes.DrawImage(gtx.Ops, im, x, image.Point{X: maxImageWidth, Y: maxImageWidth})
				// g.Pop()

				return layout.Dimensions{Size: size}
			})

			return dim
		},
	)

	// return layout.Rigid(
	// 	func(gtx models.C) models.D {
	// 		remainingSpace := (gtx.Constraints.Max.X / int(gtx.Metric.PxPerDp)) - sideBarWidth
	// 		maxWidth := remainingSpace/maxImageWidth + 1
	// 		failed := 0
	// 		for i, im := range page.images {
	// 			x := im.Bounds().Max
	// 			//fmt.Println(point, i, x)
	// 			if im == nil {
	// 				failed++
	// 				continue
	// 			}

	// 			g := op.Offset(
	// 				image.Pt(
	// 					((i * maxImageWidth) % (maxWidth * maxImageWidth)),
	// 					((i)/maxWidth)*maxImageWidth),
	// 			).Push(gtx.Ops)
	// 			shapes.DrawSquare(&gtx, shapes.SquareParams{
	// 				Size: shapes.Size{
	// 					Width:  300,
	// 					Height: 300,
	// 				},
	// 				Color: themes.MainTheme.Gray0,
	// 			})()

	// 			widget.Border{
	// 				Color:        themes.MainTheme.Gray10,
	// 				CornerRadius: unit.Dp(1),
	// 				Width:        unit.Dp(2),
	// 			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	// 				return layout.Dimensions{
	// 					Size: image.Point{
	// 						X: maxImageWidth,
	// 						Y: maxImageWidth,
	// 					},
	// 				}
	// 			})

	// 			shapes.DrawImage(gtx.Ops, im, x, image.Point{X: maxImageWidth, Y: maxImageWidth})
	// 			g.Pop()

	// 		}

	// 		d := image.Point{X: remainingSpace, Y: 2000}
	// 		return layout.Dimensions{Size: d}

	// 	},
}

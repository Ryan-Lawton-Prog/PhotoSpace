package loginPage

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace/internal/app/models"
	errorWidget "ryanlawton.art/photospace/internal/app/widgets/error"
	"ryanlawton.art/photospace/internal/app/widgets/insets"
	"ryanlawton.art/photospace/internal/app/widgets/shapes"
	"ryanlawton.art/photospace/internal/pkg/logic"
	"ryanlawton.art/photospace/internal/pkg/themes"
)

const (
	closeButtonText = "Close"
	openButtonText  = "Upload"
)

type Widgets struct {
	openButton       widget.Clickable
	closeButton      widget.Clickable
	directoryButtons []widget.Clickable
	fileButtons      []widget.Clickable
}

type Image []byte

var list = layout.List{
	Axis: layout.Vertical,
}

type Upload struct {
	widgets            Widgets
	pageQueue          chan models.PageId
	errorMessages      errorWidget.ErrorMessages
	errorMessagesQueue chan error
	position           layout.Position
	image              image.Image
	path               string
	directoryNames     []string
	fileNames          []string
	mu                 sync.Mutex
}

func NewPage(pageQueue *chan models.PageId, window *app.Window) *Upload {
	widgets := Widgets{}
	directoryNames := []string{}
	fileNames := []string{}

	path := "."

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	widgets.directoryButtons = append(widgets.directoryButtons, widget.Clickable{})
	directoryNames = append(directoryNames, "..")
	for _, e := range entries {
		if e.IsDir() {
			widgets.directoryButtons = append(widgets.directoryButtons, widget.Clickable{})
			directoryNames = append(directoryNames, e.Name())
		} else {
			widgets.fileButtons = append(widgets.fileButtons, widget.Clickable{})
			fileNames = append(fileNames, e.Name())
		}
	}

	return &Upload{
		widgets:            widgets,
		pageQueue:          *pageQueue,
		errorMessages:      errorWidget.NewErrorMessageWidget(window.Invalidate),
		errorMessagesQueue: make(chan error),
		position:           layout.Position{},
		image:              nil,
		path:               path,
		directoryNames:     directoryNames,
		fileNames:          fileNames,
	}
}

func (page *Upload) handleInput(gtx *models.C) {
	if page.widgets.closeButton.Clicked(*gtx) {
		page.pageQueue <- models.Home
	}

	if page.widgets.openButton.Clicked(*gtx) {
		// TODO: Upload file

		page.pageQueue <- models.Home
	}

	page.mu.Lock()
	for i := range len(page.widgets.directoryButtons) {
		if page.widgets.directoryButtons[i].Clicked(*gtx) {
			fmt.Println("directory: ", page.widgets.directoryButtons[i])
			// First entry is to go back a directory
			if i == 0 {
				// remove last directory in path
				path := strings.Split(page.path, "/")
				if len(path) > 1 && path[len(path)-1] != ".." {
					path = path[:len(path)-1]
				} else {
					if path[0] == "." {
						path[0] = ".."
					} else {
						path = append(path, "..")
					}
				}
				page.path = strings.Join(path, "/")
			} else {
				page.path += "/" + page.directoryNames[i]
			}

			entries, err := os.ReadDir(page.path)
			if err != nil {
				log.Fatal(err)
			}

			directoryButtons := []widget.Clickable{}
			fileButtons := []widget.Clickable{}
			directoryNames := []string{}
			fileNames := []string{}
			directoryButtons = append(directoryButtons, widget.Clickable{})
			directoryNames = append(directoryNames, "..")
			for _, e := range entries {
				if e.IsDir() {
					directoryButtons = append(directoryButtons, widget.Clickable{})
					directoryNames = append(directoryNames, e.Name())
				} else {
					fileButtons = append(fileButtons, widget.Clickable{})
					fileNames = append(fileNames, e.Name())
				}
			}
			// Lock update

			page.widgets.directoryButtons = directoryButtons
			page.widgets.fileButtons = fileButtons
			page.directoryNames = directoryNames
			page.fileNames = fileNames
			break
		}
	}

	for i := range len(page.widgets.fileButtons) {
		if page.widgets.fileButtons[i].Clicked(*gtx) {
			fmt.Println("file: ", page.widgets.fileButtons[i])
			ext := strings.Split(page.fileNames[i], ".")[len(strings.Split(page.fileNames[i], "."))-1]
			if ext == "jpg" || ext == "jpeg" || ext == "png" {
				// Upload file
				filePath := page.path + "/" + page.fileNames[i]
				file, err := os.Open(filePath)
				if err != nil {
					page.errorMessagesQueue <- err
				}
				defer file.Close()
				// Read file
				fileInfo, err := file.Stat()
				if err != nil {
					page.errorMessagesQueue <- err
				}
				fileSize := fileInfo.Size()
				blob := make([]byte, fileSize)
				_, err = file.Read(blob)
				if err != nil {
					page.errorMessagesQueue <- err
				}
				// POST file
				if logic.PostPhoto(blob, filePath) != nil {
					page.errorMessagesQueue <- fmt.Errorf("error uploading file %s", page.fileNames[i])
				}

				page.pageQueue <- models.Home
			}
		}
	}
	page.mu.Unlock()
}

func (page *Upload) StartRoutines(window *app.Window) {
	page.errorMessages.SetRefresh(window.Invalidate)
	page.errorMessages.StartRoutine()

	go func() {
		for err := range page.errorMessagesQueue {
			page.errorMessages.Add(err.Error(), time.Now().Add(time.Second*10))
		}
	}()
}

// UploadPhoto uploads a photo to the database
func (page *Upload) Layout(gtx *models.C, th *material.Theme) {
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
			// page.toolbar(&gtx, th),
			// layout.Rigid(
			// 	func(gtx models.C) layout.Dimensions {
			// 		return layout.Flex{
			// 			Axis:    layout.Horizontal,
			// 			Spacing: layout.SpaceBetween,
			// 		}.Layout(
			// 			gtx,
			// 			page.sidebar(&gtx, th),
			// 		)
			// 	},
			// ),
			page.files(&gtx, th),
			page.footer(&gtx, th),
		)

		return flex
	})

}

func (page *Upload) files(gtx *models.C, th *material.Theme) layout.FlexChild {
	return layout.Flexed(1,
		func(gtx models.C) layout.Dimensions {
			list.Layout(gtx, 1, func(gtx models.C, index int) layout.Dimensions {
				layoutList := []layout.FlexChild{}
				// lock list
				page.mu.Lock()
				for i := range len(page.widgets.directoryButtons) {
					dirMtn := material.Button(th, &page.widgets.directoryButtons[i], "📁 "+page.directoryNames[i])
					dirMtn.Inset = insets.MediumButton
					dirMtn.Background.A = 0
					dirMtn.Color = themes.MainTheme.Gray10
					layoutList = append(layoutList, layout.Rigid(dirMtn.Layout))
				}
				for i := range len(page.widgets.fileButtons) {
					ext := strings.Split(page.fileNames[i], ".")[len(strings.Split(page.fileNames[i], "."))-1]
					icon := "📄 "
					if ext == "jpg" || ext == "jpeg" || ext == "png" {
						icon = "🖼️ "
					}
					fileMtn := material.Button(th, &page.widgets.fileButtons[i], icon+page.fileNames[i])
					fileMtn.Inset = insets.MediumButton
					fileMtn.Background.A = 0
					fileMtn.Color = themes.MainTheme.Gray10
					layoutList = append(layoutList, layout.Rigid(fileMtn.Layout))
				}
				page.mu.Unlock()

				return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx, layoutList...)
			})

			return layout.Dimensions{Size: image.Point{
				X: gtx.Constraints.Max.X,
				Y: gtx.Constraints.Max.Y,
			}}
		},
	)
}

func (page *Upload) sidebar(gtx *models.C, th *material.Theme) layout.FlexChild {
	return layout.Rigid(
		func(gtx models.C) models.D {
			return layout.Dimensions{
				Size: image.Point{},
			}
		},
	)
}

func (page *Upload) footer(gtx *models.C, th *material.Theme) layout.FlexChild {
	return layout.Rigid(
		func(gtx models.C) models.D {

			// Record toolbar to get size before drawing
			macro := op.Record(gtx.Ops)
			flex := layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				// Upload Button
				layout.Rigid(
					func(gtx models.C) layout.Dimensions {
						openMtn := material.Button(th, &page.widgets.openButton, openButtonText)
						openMtn.Inset = insets.IconButton
						openMtn.Background.A = 0
						openMtn.Color = themes.MainTheme.Gray10
						uploadBtn := layout.Rigid(openMtn.Layout)

						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx, uploadBtn)
					},
				),
				layout.Rigid(
					shapes.DrawText(shapes.TextParams{
						Theme:     th,
						Text:      page.path,
						Color:     themes.MainTheme.Gray10,
						Alignment: text.Middle,
						Shadow:    true,
						Size:      shapes.H4,
					}),
				),
				layout.Rigid(
					func(gtx models.C) layout.Dimensions {
						closeMtn := material.Button(th, &page.widgets.closeButton, closeButtonText)
						closeMtn.Inset = insets.IconButton
						closeMtn.Background.A = 0
						closeMtn.Color = themes.MainTheme.Gray10
						uploadBtn := layout.Rigid(closeMtn.Layout)

						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx, uploadBtn)
					},
				),
			)
			c := macro.Stop()

			// Draw background using calculated toolbar size
			shapes.DrawSquare(&gtx, shapes.SquareParams{
				Color: themes.MainTheme.Gray2,
				Size: shapes.Size{
					Width:  gtx.Constraints.Max.X,
					Height: flex.Size.Y,
				},
			})()

			// Draw the toolbar
			c.Add(gtx.Ops)

			return flex
		},
	)
}

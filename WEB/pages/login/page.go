package loginPage

import (
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace-ui/logic"
	"ryanlawton.art/photospace-ui/models"
	"ryanlawton.art/photospace-ui/widgets/colors"
	errorWidget "ryanlawton.art/photospace-ui/widgets/error"
	"ryanlawton.art/photospace-ui/widgets/insets"
	"ryanlawton.art/photospace-ui/widgets/shapes"
)

const (
	loginTitleText   = "Please Login"
	loginButtonText  = "Login"
	signupButtonText = "Sign Up"
)

type Widgets struct {
	signupButton widget.Clickable
	loginButton  widget.Clickable
	username     widget.Editor
	password     widget.Editor
}

type Login struct {
	widgets            Widgets
	pageQueue          chan models.PageId
	errorMessages      errorWidget.ErrorMessages
	loginRequestQueue  chan logic.SignInBody
	errorMessagesQueue chan error
	loggingIn          bool
}

func NewPage(pageQueue *chan models.PageId, window *app.Window) *Login {
	return &Login{
		widgets:            Widgets{},
		pageQueue:          *pageQueue,
		errorMessages:      errorWidget.NewErrorMessageWidget(window.Invalidate),
		loginRequestQueue:  make(chan logic.SignInBody),
		errorMessagesQueue: make(chan error),
	}
}

func (page *Login) handleInput(gtx *models.C) {
	if page.widgets.loginButton.Clicked(*gtx) {
		page.loginRequestQueue <- logic.SignInBody{
			Username: page.widgets.username.Text(),
			Password: page.widgets.password.Text(),
		}
	}
}

func (page *Login) StartRoutines(window *app.Window) {
	page.errorMessages.SetRefresh(window.Invalidate)
	page.errorMessages.StartRoutine()
	go func() {
		for login := range page.loginRequestQueue {
			page.loggingIn = true

			success, err := logic.SignIn(login)
			if err != nil {
				page.errorMessagesQueue <- err
			} else if success {
				page.pageQueue <- models.Home
			}

			page.loggingIn = false
		}
	}()

	go func() {
		for err := range page.errorMessagesQueue {
			page.errorMessages.Add(err.Error(), time.Now().Add(time.Second*10))
		}
	}()
}

// UploadPhoto uploads a photo to the database
func (page *Login) Layout(gtx *models.C, th *material.Theme) {
	page.handleInput(gtx)

	layout.Flex{
		// Vertical alignment, from top to bottom
		Axis: layout.Vertical,
		// Empty space is left at the start, i.e. at the top
		Spacing: layout.SpaceEnd,
	}.Layout(*gtx,
		// TITLE
		layout.Rigid(
			func(gtx models.C) layout.Dimensions {
				return shapes.DrawText(&gtx, shapes.Params{
					Theme:     th,
					Text:      loginTitleText,
					Color:     colors.MainTheme.Gray10,
					Alignment: text.Middle,
					Shadow:    true,
				})
			},
		),
		layout.Rigid(
			page.errorMessages.Layout(gtx, th),
		),
		// USERNAME TEXT BOX
		layout.Rigid(
			func(gtx models.C) models.D {
				return shapes.DrawTextBox(&gtx, &page.widgets.username, shapes.Params{
					Theme: th,
					Text:  "username",
				})
			},
		),
		// PASSWORD TEXT BOX
		layout.Rigid(
			func(gtx models.C) models.D {
				return shapes.DrawTextBox(&gtx, &page.widgets.password, shapes.Params{
					Theme: th,
					Text:  "password",
				})
			},
		),
		// SIGNUP AND SUBMIT BUTTONS
		layout.Rigid(
			func(gtx models.C) layout.Dimensions {
				// ONE: First define margins around the button using layout.Inset ...
				margins := layout.Inset{
					Top:    unit.Dp(25),
					Bottom: unit.Dp(25),
					Right:  unit.Dp(35),
					Left:   unit.Dp(35),
				}
				// TWO: ... then we lay out those margins ...
				return margins.Layout(gtx,
					// THREE: ... and finally within the margins, we define and lay out the button
					func(gtx models.C) models.D {
						loginMat := material.Button(th, &page.widgets.loginButton, loginButtonText)
						loginMat.Inset = insets.LargeButton
						loginMat.Background = colors.MainTheme.Primary5
						loginMat.Color = colors.White
						signupMat := material.Button(th, &page.widgets.signupButton, signupButtonText)
						signupMat.Inset = insets.LargeButton
						signupMat.Background = colors.MainTheme.Primary5
						signupMat.Color = colors.White
						loginBtn := layout.Rigid(loginMat.Layout)
						signupBtn := layout.Rigid(signupMat.Layout)

						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly}.Layout(gtx, loginBtn, signupBtn)
					},
				)
			},
		),
		// ... then one to hold an empty spacer
		layout.Rigid(
			// The height of the spacer is 25 Device independent pixels
			layout.Spacer{Height: unit.Dp(25)}.Layout,
		),
	)
}

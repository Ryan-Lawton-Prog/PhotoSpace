package page

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace/internal/app/models"
	homePage "ryanlawton.art/photospace/internal/app/pages/home"
	loginPage "ryanlawton.art/photospace/internal/app/pages/login"
)

type IPage interface {
	Layout(gtx *models.C, th *material.Theme)
	StartRoutines(win *app.Window)
}

func GetPageFactory(page models.PageId, pageQueue *chan models.PageId, window *app.Window) (IPage, error) {
	switch page {
	case models.Login:
		return loginPage.NewPage(pageQueue, window), nil
	case models.Home:
		return homePage.NewPage(pageQueue, window), nil
	}

	return nil, fmt.Errorf("wrong page id passed")
}

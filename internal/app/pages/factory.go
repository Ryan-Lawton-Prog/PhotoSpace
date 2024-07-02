package page

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace/internal/app/models"
	loginPage "ryanlawton.art/photospace/internal/app/pages/login"
	snakePage "ryanlawton.art/photospace/internal/app/pages/snake"
)

type IPage interface {
	Layout(gtx *models.C, th *material.Theme)
	StartRoutines(win *app.Window)
}

func GetPageFactory(page models.PageId, pageQueue *chan models.PageId, window *app.Window) (IPage, error) {
	switch page {
	case models.Login:
		return loginPage.NewPage(pageQueue, window), nil
	case models.Snake:
		return snakePage.NewPage(pageQueue), nil
	}

	return nil, fmt.Errorf("wrong page id passed")
}

package router

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace-ui/models"
	page "ryanlawton.art/photospace-ui/pages"
)

type Router struct {
	Window    *app.Window
	Page      page.IPage
	PageQueue chan models.PageId
	Theme     material.Theme
}

func NewRouter(window *app.Window, initPage models.PageId) Router {
	r := Router{
		Window:    window,
		PageQueue: make(chan models.PageId),
		Theme:     *material.NewTheme(),
	}

	r.Page, _ = page.GetPageFactory(initPage, &r.PageQueue, window)
	return r
}

func (router *Router) Loop() {
	th := material.NewTheme()
	var ops op.Ops
	router.Route()
	router.Page.StartRoutines(router.Window)

	for {
		switch e := router.Window.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			router.Page.Layout(&gtx, th)

			e.Frame(gtx.Ops)
		}
	}
}

func (router *Router) Route() {
	go func() {
		for pageId := range router.PageQueue {
			fmt.Println("page: ", pageId)
			router.Page, _ = page.GetPageFactory(pageId, &router.PageQueue, router.Window)
			router.Page.StartRoutines(router.Window)
			router.Window.Invalidate()
		}
	}()
}

package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"ryanlawton.art/photospace/internal/app/models"
	"ryanlawton.art/photospace/internal/app/router"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Title("PhotoSpace"))
		window.Option(app.Size(unit.Dp(1920), unit.Dp(1080)))
		//window.Option(app.Maximized.Option())
		r := router.NewRouter(window, models.Login)
		r.Loop()
		os.Exit(0)
	}()

	app.Main()
}

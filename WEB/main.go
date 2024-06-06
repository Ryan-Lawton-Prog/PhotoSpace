package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"ryanlawton.art/photospace-ui/models"
	"ryanlawton.art/photospace-ui/router"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Title("PhotoSpace"))
		window.Option(app.Size(unit.Dp(400), unit.Dp(600)))
		r := router.NewRouter(window, models.Login)
		r.Loop()
		os.Exit(0)
	}()

	app.Main()
}

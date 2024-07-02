package models

import (
	"gioui.org/layout"
)

type C = layout.Context
type D = layout.Dimensions

type Point struct {
	X, Y float32
}

type Rectangle struct {
	Min, Max Point
}

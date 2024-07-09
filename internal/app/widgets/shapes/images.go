package shapes

import (
	"fmt"
	"image"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"
)

func DrawImage(ops *op.Ops, img image.Image, bounds image.Point, maxSize image.Point) {
	imageOp := paint.NewImageOp(img)
	imageOp.Filter = paint.FilterNearest
	imageOp.Add(ops)
	scalerX := 1 / (float32(bounds.X) / float32(maxSize.X))
	scalerY := 1 / (float32(bounds.Y) / float32(maxSize.Y))
	bestScaler := min(scalerX, scalerY)
	fmt.Println(((float32(maxSize.Y) - (float32(bounds.Y) * bestScaler)) / 2), (float32(bounds.Y) * bestScaler))
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, ((float32(maxSize.Y)-(float32(bounds.Y)*bestScaler))/2)), f32.Pt(bestScaler, bestScaler))).Add(ops)
	paint.PaintOp{}.Add(ops)
}

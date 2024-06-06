package snakePage

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace-ui/models"
	"ryanlawton.art/photospace-ui/widgets/shapes"
)

type Pos [2]int
type Col []bool
type Board []Col

type Widgets struct {
	startButton widget.Clickable
	backButton  widget.Clickable
}
type Snake struct {
	widgets       Widgets
	pageQueue     chan models.PageId
	board         Board
	food          Pos
	snake         []Pos
	snakeLength   int
	lastDirection cardinal
	gameOver      bool
	gamePlaying   bool
}

const BoardWidth, BoardHeight = 8, 6

type Tile struct {
	Color color.NRGBA
	Size  int
}

var squareColor = color.NRGBA{A: 0xFF}
var snakeColor = color.NRGBA{G: 0x80, A: 0xFF}
var foodColor = color.NRGBA{R: 0x80, A: 0xFF}

type cardinal int

const (
	Left cardinal = iota
	Right
	Up
	Down
)

var keyToCardinal = map[key.Name]cardinal{
	key.NameLeftArrow:  Left,
	key.NameRightArrow: Right,
	key.NameUpArrow:    Up,
	key.NameDownArrow:  Down,
}

func NewPage(pageQueue *chan models.PageId) *Snake {
	board := make(Board, BoardWidth)
	for i := range board {
		board[i] = make(Col, BoardHeight)
	}
	board[0][0] = true

	return &Snake{
		widgets:       Widgets{},
		pageQueue:     *pageQueue,
		board:         board,
		food:          randomPos(BoardWidth-1, BoardHeight-1),
		snake:         []Pos{{0, 0}},
		snakeLength:   1,
		gameOver:      false,
		lastDirection: Right,
	}
}

// UploadPhoto uploads a photo to the database
func (screen *Snake) Layout(gtx *models.C, th *material.Theme) {
	if screen.widgets.backButton.Clicked(*gtx) {
		screen.pageQueue <- models.Login
	}

	if screen.widgets.startButton.Clicked(*gtx) {
		screen.gamePlaying = !screen.gamePlaying
	}

	screen.changeDirection(gtx.Event(key.Filter{
		Name: key.NameRightArrow,
	}))
	screen.changeDirection(gtx.Event(key.Filter{
		Name: key.NameLeftArrow,
	}))
	screen.changeDirection(gtx.Event(key.Filter{
		Name: key.NameUpArrow,
	}))
	screen.changeDirection(gtx.Event(key.Filter{
		Name: key.NameDownArrow,
	}))

	// Define an large label with an appropriate text:
	title := material.H1(th, "Hello, Snake")

	layout.Flex{
		// Vertical alignment, from top to bottom
		Axis: layout.Vertical,

		// Empty space is left at the start, i.e. at the top
		Spacing: layout.SpaceStart,

		// Horizontal alignment
		Alignment: layout.Alignment(layout.Middle),
	}.Layout(*gtx,
		layout.Rigid(
			func(gtx models.C) models.D {
				if screen.gameOver {
					return shapes.DrawText(&gtx, shapes.Params{
						Theme:     th,
						Text:      "GAME OVER",
						Color:     color.NRGBA{R: 127, G: 0, B: 0, A: 255},
						Alignment: text.Middle,
					})
				}

				return layout.Dimensions{}
			},
		),
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
						btn := material.Button(th, &screen.widgets.backButton, "Egg")
						return btn.Layout(gtx)
					},
				)
			},
		),
		layout.Rigid(
			func(gtx models.C) models.D {

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
						if !screen.gameOver {
							for i := range screen.board {
								for j := range screen.board[i] {
									tile := screen.getTile(Pos{i, j})
									shapes.OffsetDraw(
										&gtx,
										shapes.Offset{120*i + (120 - tile.Size/2), 120*j + (120 - tile.Size/2)},
										shapes.DrawSquare,
										shapes.Params{
											Color: tile.Color,
											Size:  shapes.Size{Width: tile.Size, Height: tile.Size},
										},
									)
								}
							}

							d := image.Point{X: len(screen.board) * 120, Y: len(screen.board[0]) * 120}
							return layout.Dimensions{Size: d}
						}

						return layout.Dimensions{}

					},
				)
			},
		),
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
						text := "Start"
						if screen.gamePlaying {
							text = "Pause"
						}
						btn := material.Button(th, &screen.widgets.startButton, text)
						return btn.Layout(gtx)
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

	// Change the color of the label.
	maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	title.Color = maroon

	// Change the position of the label.
	title.Alignment = text.Middle

	// Draw the label to the graphics context.
	title.Layout(*gtx)
}

func (screen *Snake) StartRoutines(window *app.Window) {
	go func() {
		for {
			if screen.gamePlaying {
				screen.moveSnake(screen.lastDirection)
				window.Invalidate()
				mul := 1000 - (10 * (screen.snakeLength * 2))
				time.Sleep(time.Millisecond * time.Duration(mul))
			}
		}
	}()
}

func (screen *Snake) changeDirection(event event.Event, ok bool) {
	if !ok {
		return
	}
	screen.lastDirection = keyToCardinal[event.(key.Event).Name]
}

func (screen *Snake) moveSnake(direction cardinal) {
	var newHead = screen.snake[screen.snakeLength-1]
	switch direction {
	case Up:
		newHead[1] -= 1
	case Down:
		newHead[1] += 1
	case Left:
		newHead[0] -= 1
	case Right:
		newHead[0] += 1
	}

	if newHead[0] < 0 ||
		newHead[1] < 0 ||
		newHead[0] >= len(screen.board) ||
		newHead[1] >= len(screen.board[0]) ||
		screen.board[newHead[0]][newHead[1]] {
		screen.gameOver = true
		return
	}

	screen.lastDirection = direction

	screen.snake = append(screen.snake, newHead)
	screen.board[newHead[0]][newHead[1]] = true

	if newHead != screen.food {
		screen.board[screen.snake[0][0]][screen.snake[0][1]] = false
		screen.snake = screen.snake[1:]
		return
	}

	screen.snakeLength++
	screen.spawnFood()
}

func (screen *Snake) spawnFood() {
	for screen.board[screen.food[0]][screen.food[1]] {
		screen.food = randomPos(len(screen.board), len(screen.board[0]))
	}
}

func randomPos(xMax, yMax int) Pos {
	return Pos{rand.Int() % xMax, rand.Int() % yMax}
}

func (screen *Snake) getTile(pos Pos) Tile {
	if screen.food == pos {
		return Tile{
			Color: foodColor,
			Size:  50,
		}
	}

	if screen.board[pos[0]][pos[1]] {
		return Tile{
			Color: snakeColor,
			Size:  90,
		}
	}

	return Tile{
		Color: squareColor,
		Size:  100,
	}
}

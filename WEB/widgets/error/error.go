package errorWidget

import (
	"fmt"
	"slices"
	"sort"
	"sync"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	"ryanlawton.art/photospace-ui/models"
	"ryanlawton.art/photospace-ui/widgets/colors"
	"ryanlawton.art/photospace-ui/widgets/shapes"
)

type ErrorMessage struct {
	message  string
	duration time.Time
}

type ErrorMessages struct {
	mu       sync.Mutex
	messages []ErrorMessage
	refresh  func()
}

func NewErrorMessageWidget(refresh func()) ErrorMessages {
	return ErrorMessages{
		refresh: refresh,
	}
}

func (messages *ErrorMessages) SetRefresh(r func()) {
	messages.refresh = r
}

func (messages *ErrorMessages) Add(message string, duration time.Time) {
	fmt.Println("Adding error: ", message)
	messages.mu.Lock()
	defer messages.mu.Unlock()
	index := sort.Search(len(messages.messages), func(i int) bool {
		return messages.messages[i].duration.Unix() > duration.Unix()
	})
	messages.messages = slices.Insert(messages.messages, index, ErrorMessage{
		message:  message,
		duration: duration,
	})

	messages.refresh()
}

func (messages *ErrorMessages) StartRoutine() {
	go func() {
		// Remove old messages
		for {
			messages.mu.Lock()
			for len(messages.messages) > 0 && messages.messages[0].duration.Unix() <= time.Now().Unix() {
				fmt.Println("Removing message", messages.messages[0])
				messages.messages = messages.messages[1:]
			}
			messages.mu.Unlock()
			messages.refresh()
			time.Sleep(time.Second)
		}
	}()
}

func (messages *ErrorMessages) GetErrors() []string {
	messages.mu.Lock()
	defer messages.mu.Unlock()
	var widgets []string
	for _, message := range messages.messages {
		widgets = append(widgets, message.message)
	}

	return widgets
}

func (messages *ErrorMessages) Layout(gtx *models.C, th *material.Theme) func(gtx models.C) layout.Dimensions {
	return func(gtx models.C) layout.Dimensions {
		var d layout.Dimensions
		for _, message := range messages.GetErrors() {
			textDi := shapes.DrawText(&gtx, shapes.Params{
				Theme:     th,
				Text:      message,
				Color:     colors.Red,
				Alignment: text.Middle,
			})

			d.Size = d.Size.Add(textDi.Size)
		}
		return d
	}
}

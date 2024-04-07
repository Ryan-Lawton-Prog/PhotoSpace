package models

import "image"

type Photo struct {
	ID        string
	UserID    string
	AlbumID   string
	Photo     image.Image
	Extension string
	Title     string
}

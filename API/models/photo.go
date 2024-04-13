package models

type Photo struct {
	ID       string
	UserID   string
	AlbumID  string
	Photo    []byte
	Filename string
}

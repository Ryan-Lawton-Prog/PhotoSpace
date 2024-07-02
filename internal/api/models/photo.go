package models

type PhotoMetadata struct {
	ID        string
	UserID    string
	AlbumID   string
	Filename  string
	BucketURL string
}

type PhotoBlob []byte

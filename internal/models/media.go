package models

import "io"

type MediaFile struct {
	Filename string
	File     io.ReadCloser
}

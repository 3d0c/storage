package saver

import (
	"io"
)

// Saver interface
type Saver interface {
	Save(r io.Reader) error
}

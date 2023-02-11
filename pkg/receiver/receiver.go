package receiver

import (
	"io"
)

// Receiver interface
type Receiver interface {
	Recv(dst io.Writer) error
}

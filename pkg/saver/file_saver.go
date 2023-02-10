package saver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileSaver struct
type FileSaver struct {
	objectID string
	dir      string
}

// NewFileSaver constructor
func NewFileSaver(objectID, dir string) (*FileSaver, error) {
	var (
		fs = &FileSaver{
			objectID: objectID,
			dir:      dir,
		}
		err error
	)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("error creating storage directory '%s' - %s", dir, err)
		}
	}

	return fs, nil
}

// Save implements Saver interface
func (fs *FileSaver) Save(src io.Reader) error {
	var (
		dst *os.File
		err error
	)

	filePath := filepath.Join(fs.dir, fs.objectID)

	if dst, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return fmt.Errorf("error creating file '%s' - %s", filePath, err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("error saving file '%s' - %s", filePath, err)
	}

	return nil
}

package saver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/3d0c/storage/pkg/utils"
)

// FileSaver struct
type FileSaver struct {
	objectID string
	dir      string
}

// NewFileSaver constructor
func NewFileSaver(objectID, root string) (*FileSaver, error) {
	var (
		fs = &FileSaver{
			objectID: objectID,
		}
		err error
	)

	fs.dir = filepath.Dir(utils.BuildFilePath(root, objectID))

	if _, err = os.Stat(fs.dir); os.IsNotExist(err) {
		if err = os.MkdirAll(fs.dir, 0755); err != nil {
			return nil, fmt.Errorf("error creating storage directory '%s' - %s", fs.dir, err)
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

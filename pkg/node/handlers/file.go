package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/saver"
)

// File is a file "namespace"
type File struct{}

// FileHandler "file namespace" constructor
func FileHandler() *File {
	return &File{}
}

// Put handler implementation
func (*File) Put(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		objectID = chi.URLParam(r, "ID")
		sv       saver.Saver
		err      error
	)

	if sv, err = saver.NewFileSaver(objectID, config.Node().StorageDir); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing FileSaver - %s", err)
	}

	defer r.Body.Close()

	if err = sv.Save(r.Body); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error saving data - %s", err)
	}

	return nil, http.StatusOK, nil
}

// Get handler implementation
func (*File) Get(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		objectID = chi.URLParam(r, "ID")
		src      *os.File
		err      error
	)

	filePath := filepath.Join(config.Node().StorageDir, objectID)

	if src, err = os.OpenFile(filePath, os.O_RDONLY, 0644); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error opening file '%s' - %s", filePath, err)
	}
	defer src.Close()

	if _, err = io.Copy(w, src); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error copying file '%s' - %s", filePath, err)
	}

	return nil, http.StatusOK, nil
}

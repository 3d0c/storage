package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/saver"
	"github.com/3d0c/storage/pkg/utils"
)

// File structure. Namespace and config storage
type File struct {
	cfg config.NodeConfig
}

// FileHandler file struct constructor
func FileHandler(c config.NodeConfig) *File {
	return &File{
		cfg: c,
	}
}

// Put handler implementation
func (f *File) Put(w http.ResponseWriter, r *http.Request) (int, error) {
	var (
		objectID = chi.URLParam(r, "ID")
		sv       saver.Saver
		err      error
	)

	if sv, err = saver.NewFileSaver(objectID, f.cfg.StorageDir); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error initializing FileSaver - %s", err)
	}

	if err = sv.Save(r.Body); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error saving data - %s", err)
	}

	return http.StatusOK, nil
}

// Get handler implementation
func (f *File) Get(w http.ResponseWriter, r *http.Request) (int, error) {
	var (
		objectID = chi.URLParam(r, "ID")
		src      *os.File
		err      error
	)

	filePath := utils.BuildFilePath(f.cfg.StorageDir, objectID)

	if src, err = os.OpenFile(filePath, os.O_RDONLY, 0644); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error opening file '%s' - %s", filePath, err)
	}
	defer src.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	if _, err = io.Copy(w, src); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error copying file '%s' - %s", filePath, err)
	}

	return http.StatusOK, nil
}

// Delete handler implementation. Used as a rollback endpoint
func (f *File) Delete(_ http.ResponseWriter, r *http.Request) (int, error) {
	var (
		objectID = chi.URLParam(r, "ID")
		err      error
	)

	filePath := utils.BuildFilePath(f.cfg.StorageDir, objectID)

	if err = os.Remove(filePath); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error removing file '%s' - %s", filePath, err)
	}

	return http.StatusOK, nil
}

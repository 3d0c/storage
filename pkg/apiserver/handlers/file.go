package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/database/models"
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
		nodes    = config.Proxy().Nodes
		sv       saver.Saver
		pm       *models.Parts
		size     int
		err      error
	)

	if pm, err = models.NewPartsModel(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing Parts model - %s", err)
	}

	if _, err = pm.FindNodes(objectID); err == models.ErrNotFound {
		return nil, http.StatusForbidden, fmt.Errorf("object '%s' already exist", objectID)
	}

	if size, err = strconv.Atoi(r.Header.Get("Content-Length")); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting object size - %s", err)
	}

	if sv, err = saver.NewHTTPSaver(objectID, size, nodes); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing HTTPSaver - %s", err)
	}

	defer r.Body.Close()

	if err = sv.Save(r.Body); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error saving data - %s", err)
	}

	return nil, http.StatusOK, nil
}

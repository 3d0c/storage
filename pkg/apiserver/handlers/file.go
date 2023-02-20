package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/database/models"
	"github.com/3d0c/storage/pkg/receiver"
	"github.com/3d0c/storage/pkg/saver"
)

// File structure. Namespace and config storage
type File struct {
	cfg config.ProxyConfig
}

// FileHandler file struct constructor
func FileHandler(c config.ProxyConfig) *File {
	return &File{
		cfg: c,
	}
}

// Put handler implementation
func (f *File) Put(_ http.ResponseWriter, r *http.Request) (int, error) {
	var (
		objectID = chi.URLParam(r, "ID")
		sv       saver.Saver
		om       *models.Object
		size     int
		err      error
	)

	if om, err = models.NewObjectModel(f.cfg.Database); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error initializing Object model - %s", err)
	}

	if _, err = om.Find(objectID); err == nil {
		return http.StatusForbidden, fmt.Errorf("object %s already exist", objectID)
	}

	if size, err = strconv.Atoi(r.Header.Get("Content-Length")); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error getting object size - %s", err)
	}

	if sv, err = saver.NewHTTPSaver(objectID, size, f.cfg); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error initializing HTTPSaver - %s", err)
	}

	if err = sv.Save(r.Body); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error saving data - %s", err)
	}

	if err = om.Add(objectID, size); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error adding object %s - %s", objectID, err)
	}

	return http.StatusOK, nil
}

// Get handler implementation
func (f *File) Get(w http.ResponseWriter, r *http.Request) (int, error) {
	var (
		objectID = chi.URLParam(r, "ID")
		pm       *models.Parts
		om       *models.Object
		object   models.Object
		rv       receiver.Receiver
		nodes    []int
		err      error
	)

	if om, err = models.NewObjectModel(f.cfg.Database); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error initializing Object model - %s", err)
	}

	// Get Object by it's id. It's needed to get it's length to setup Content-Lenght
	if object, err = om.Find(objectID); err != nil {
		if err == models.ErrNotFound {
			return http.StatusNotFound, fmt.Errorf("object '%s' not found", objectID)
		}
		return http.StatusInternalServerError, fmt.Errorf("error finding object '%s' - %s", objectID, err)
	}

	if pm, err = models.NewPartsModel(f.cfg.Database); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error initializing Parts model - %s", err)
	}

	// Get nodes where parts are stored
	if nodes, err = pm.FindNodes(objectID); err != nil {
		if err == models.ErrNotFound {
			return http.StatusNotFound, fmt.Errorf("nodes for object '%s' not found", objectID)
		}
		return http.StatusInternalServerError, fmt.Errorf("error getting nodes for object '%s' - %s", objectID, err)
	}

	// Receiver here stands for getting parts(chunks) from remote nodes
	if rv, err = receiver.NewHTTPReceiver(objectID, nodes, f.cfg); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error initializing receiver - %s", err)
	}

	// Set Content-Length first
	w.Header().Set("Content-Length", strconv.Itoa(object.Size))

	if err = rv.Recv(w); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error receiving and sending object '%s' - %s", objectID, err)
	}

	return http.StatusOK, nil
}

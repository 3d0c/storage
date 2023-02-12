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
		om       *models.Object
		size     int
		err      error
	)

	if om, err = models.NewObjectModel(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing Object model - %s", err)
	}

	if _, err = om.Find(objectID); err == nil {
		return nil, http.StatusForbidden, fmt.Errorf("object %s already exist", objectID)
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

	if err = om.Add(objectID, size); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error adding object %s - %s", objectID, err)
	}

	return nil, http.StatusOK, nil
}

// Get handler implementation
func (*File) Get(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		objectID = chi.URLParam(r, "ID")
		pm       *models.Parts
		om       *models.Object
		object   models.Object
		rv       receiver.Receiver
		nodes    []int
		err      error
	)

	if om, err = models.NewObjectModel(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing Object model - %s", err)
	}

	// Get Object by it's id. It's needed to get it's length to setup Content-Lenght
	if object, err = om.Find(objectID); err != nil {
		if err == models.ErrNotFound {
			return nil, http.StatusNotFound, fmt.Errorf("object '%s' not found", objectID)
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("error finding object '%s' - %s", objectID, err)
	}

	if pm, err = models.NewPartsModel(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing Parts model - %s", err)
	}

	// Get nodes where parts are stored
	if nodes, err = pm.FindNodes(objectID); err != nil {
		if err == models.ErrNotFound {
			return nil, http.StatusNotFound, fmt.Errorf("nodes for object '%s' not found", objectID)
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting nodes for object '%s' - %s", objectID, err)
	}

	// Receiver here stands for getting parts(chunks) from remote nodes
	if rv, err = receiver.NewHTTPReceiver(objectID, nodes); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing receiver - %s", err)
	}

	// Set Content-Length first
	w.Header().Set("Content-Length", strconv.Itoa(object.Size))

	if err = rv.Recv(w); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error receiving and sending object '%s' - %s", objectID, err)
	}

	return nil, http.StatusOK, nil
}

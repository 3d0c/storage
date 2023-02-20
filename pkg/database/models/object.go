package models

import (
	"database/sql"
	"fmt"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/database"
)

// Object Model
type Object struct {
	cfg      config.Database
	ObjectID string
	Size     int
}

// NewObjectModel Object constructor
func NewObjectModel(cfg config.Database) (*Object, error) {
	return &Object{
		cfg: cfg,
	}, nil
}

// Find finds object by id
func (o *Object) Find(objectID string) (Object, error) {
	var (
		result Object
		err    error
		stmt   string = "SELECT object_id, size FROM objects WHERE object_id = ?"
	)

	if err = database.Instance(o.cfg).QueryRow(stmt, objectID).Scan(&result.ObjectID, &result.Size); err != nil {
		if err == sql.ErrNoRows {
			return Object{}, ErrNotFound
		}
	}

	return result, nil
}

// Add an object
func (o *Object) Add(objectID string, size int) error {
	var (
		err  error
		stmt string = "INSERT INTO objects VALUES (?,?)"
	)

	if _, err = database.Instance(o.cfg).Exec(stmt, objectID, size); err != nil {
		return fmt.Errorf("error inserting object - %s", err)
	}

	return nil
}

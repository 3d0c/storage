package models

import (
	"database/sql"
	"fmt"

	"github.com/3d0c/storage/pkg/database"
)

// Object Model
type Object struct {
	ObjectID string
	Size     int
}

// NewObjectModel Object constructor
func NewObjectModel() (*Object, error) {
	return &Object{}, nil
}

// Find finds object by id
func (*Object) Find(objectID string) (Object, error) {
	var (
		result Object
		err    error
		stmt   string = "SELECT object_id, size FROM objects WHERE object_id = ?"
	)

	if err = database.Instance().QueryRow(stmt, objectID).Scan(&result.ObjectID, &result.Size); err != nil {
		if err == sql.ErrNoRows {
			return Object{}, ErrNotFound
		}
	}

	return result, nil
}

// Add an object
func (*Object) Add(objectID string, size int) error {
	var (
		err  error
		stmt string = "INSERT INTO objects VALUES (?,?)"
	)

	if _, err = database.Instance().Exec(stmt, objectID, size); err != nil {
		return fmt.Errorf("error inserting object - %s", err)
	}

	return nil
}

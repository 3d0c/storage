package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/3d0c/storage/pkg/database"
)

// ErrNotFound return this type to not to expose sql package into handlers
var ErrNotFound = errors.New("nothing found")

// Parts Model
type Parts struct{}

// NewPartsModel Parts constructor
func NewPartsModel() (*Parts, error) {
	return &Parts{}, nil
}

// FindNodes for object
func (*Parts) FindNodes(objectID string) ([]int, error) {
	var (
		result = make([]int, 0)
		nodeID int
		rows   *sql.Rows
		err    error
		stmt   string = "SELECT node_id FROM parts WHERE object_id = ? ORDER BY part_id ASC"
	)

	if rows, err = database.Instance().Query(stmt, objectID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&nodeID); err != nil {
			return nil, fmt.Errorf("error scanning row - %s", err)
		}

		result = append(result, nodeID)
	}

	return result, nil
}

// Add part
func (*Parts) Add(objectID string, nodeID int, partID int) error {
	var (
		err  error
		stmt string = "INSERT INTO parts VALUES (?,?,?)"
	)

	if _, err = database.Instance().Exec(stmt, objectID, nodeID, partID); err != nil {
		return fmt.Errorf("error inserting part - %s", err)
	}

	return nil
}

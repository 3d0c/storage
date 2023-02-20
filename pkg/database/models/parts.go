package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/database"
)

// ErrNotFound return this type to not to expose sql package into handlers
var ErrNotFound = errors.New("nothing found")

// Parts Model
type Parts struct {
	cfg config.Database
	tx  *sql.Tx
}

// NewPartsModel Parts constructor
func NewPartsModel(cfg config.Database) (*Parts, error) {
	var (
		parts = &Parts{
			cfg: cfg,
		}
		err error
	)

	if parts.tx, err = database.Instance(cfg).Begin(); err != nil {
		return nil, fmt.Errorf("error starting transation - %s", err)
	}

	return parts, nil
}

// FindNodes for object
func (p *Parts) FindNodes(objectID string) ([]int, error) {
	var (
		result = make([]int, 0)
		nodeID int
		rows   *sql.Rows
		err    error
		stmt   string = "SELECT node_id FROM parts WHERE object_id = ? ORDER BY part_id ASC"
	)
	defer func() {
		if err := p.Commit(); err != nil {
			fmt.Printf("error commiting transaction - %s", err)
		}
	}()

	if rows, err = database.Instance(p.cfg).Query(stmt, objectID); err != nil {
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
// No default TX action. TX should be closed on requesting side
func (p *Parts) Add(objectID string, nodeID int, partID int) error {
	var (
		err  error
		stmt string = "INSERT INTO parts VALUES (?,?,?)"
	)

	if _, err = database.Instance(p.cfg).Exec(stmt, objectID, nodeID, partID); err != nil {
		return fmt.Errorf("error inserting part - %s", err)
	}

	return nil
}

// Commit wrapper
func (p *Parts) Commit() error {
	if err := p.tx.Commit(); err != nil {
		return fmt.Errorf("error commiting transaction - %s", err)
	}
	return nil
}

// Rollback wrapper
func (p *Parts) Rollback() error {
	if err := p.tx.Rollback(); err != nil {
		return fmt.Errorf("error rollbacking transaction - %s", err)
	}
	return nil
}

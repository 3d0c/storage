package config

import (
	"time"

	"github.com/3d0c/storage/pkg/utils"
)

// Server section
type Server struct {
	Address     string
	ReadTimeout time.Duration
	WriteTemout time.Duration
}

// Logger section
type Logger struct {
	Level       string
	AddCaller   bool
	OutputPaths []string
}

// Database section
type Database struct {
	DSN     string
	Dialect string
}

// Saver section
type Saver struct {
	StorageDir string
}

// Nodes section
type Nodes []string

// Pick gets random node. Return node address and it's id
func (n Nodes) Pick() (string, int) {
	id := utils.RandomSeed(0, len(n)-1)
	return n[id], id
}

// Get node by id
func (n Nodes) Get(id int) string {
	if id > len(n)-1 {
		return ""
	}
	return n[id]
}

package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/log"
	"github.com/3d0c/storage/pkg/utils"
)

var (
	expectedNodes    []int
	expectedObjectID string
)

func TestFindBy(t *testing.T) {
	var (
		pm    *Parts
		nodes []int
		err   error
	)

	pm, err = NewPartsModel()
	assert.Nil(t, err)

	nodes, err = pm.FindNodes(expectedObjectID)
	assert.Nil(t, err)
	assert.Equal(t, expectedNodes, nodes)
}

func prepareDatabase() error {
	var (
		rndSample = utils.RandomSeed(0, 10)
		nodeIDs   []int
		objectID  string
		pm        *Parts
		err       error
	)

	if pm, err = NewPartsModel(); err != nil {
		return fmt.Errorf("error initializing Parts model - %s", err)
	}

	// Gerenarate 10 objects
	for i := 0; i < 10; i++ {
		objectID = utils.RandomString(32)
		nodeIDs = nil

		// with 5 random nodes of 10 available
		for j := 0; j < 5; j++ {
			nodeID := utils.RandomSeed(0, 9)
			if err = pm.Add(objectID, nodeID, j); err != nil {
				panic(err)
			}
			nodeIDs = append(nodeIDs, nodeID)
		}

		if i == rndSample {
			expectedObjectID = objectID
			expectedNodes = nodeIDs
		}
	}

	return nil
}

func TestMain(m *testing.M) {
	const (
		dbFileName = "/tmp/models_test.db"
	)
	var (
		err error
	)

	log.InitLogger(config.Logger{
		Level:     "debug",
		AddCaller: true,
	})

	config.Proxy().Database.DSN = dbFileName
	config.Proxy().Database.Dialect = "sqlite3"

	if err = prepareDatabase(); err != nil {
		fmt.Printf("Error preparing testing environment - %s\n", err)
		os.Exit(-1)
	}

	exitval := m.Run()

	if err := os.Remove(dbFileName); err != nil {
		fmt.Printf("ERROR: removing testing database, error - %s\n", err)
	}

	os.Exit(exitval)
}

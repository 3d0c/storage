package config

import (
	"sync"
)

var (
	nodeInstance *NodeConfig
	nodeOnce     sync.Once
)

// NodeConfig config
type NodeConfig struct {
	Server
	Logger
	Saver
}

// Node instance
func Node() *NodeConfig {
	nodeOnce.Do(func() {
		nodeInstance = new(NodeConfig)
	})

	return nodeInstance
}

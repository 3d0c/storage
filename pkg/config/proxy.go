package config

import (
	"sync"
)

var (
	proxyInstance *ProxyConfig
	proxyOnce     sync.Once
)

// ProxyConfig config
type ProxyConfig struct {
	Server
	Logger
	Database
	Nodes
}

// Proxy instance
func Proxy() *ProxyConfig {
	proxyOnce.Do(func() {
		proxyInstance = new(ProxyConfig)
	})

	return proxyInstance
}

.PHONY: build

.DEFAULT_GOAL=build

version:
	echo $(FULL_VERSION)

build/proxy:
	go build -o build/proxy proxy_main.go

build/node:
	go build -o build/node node_main.go

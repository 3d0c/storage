#!/bin/sh

NODE_LISTEN=":9000" NODE_STORAGE="/tmp/node-storage-0" ./node_main run --config=configs/dev/node.yaml&
NODE_LISTEN=":9001" NODE_STORAGE="/tmp/node-storage-1" ./node_main run --config=configs/dev/node.yaml&
NODE_LISTEN=":9002" NODE_STORAGE="/tmp/node-storage-2" ./node_main run --config=configs/dev/node.yaml&
NODE_LISTEN=":9003" NODE_STORAGE="/tmp/node-storage-3" ./node_main run --config=configs/dev/node.yaml&
NODE_LISTEN=":9004" NODE_STORAGE="/tmp/node-storage-4" ./node_main run --config=configs/dev/node.yaml&
NODE_LISTEN=":9005" NODE_STORAGE="/tmp/node-storage-5" ./node_main run --config=configs/dev/node.yaml&
NODE_LISTEN=":9006" NODE_STORAGE="/tmp/node-storage-6" ./node_main run --config=configs/dev/node.yaml&
NODE_LISTEN=":9007" NODE_STORAGE="/tmp/node-storage-7" ./node_main run --config=configs/dev/node.yaml&
NODE_LISTEN=":9008" NODE_STORAGE="/tmp/node-storage-8" ./node_main run --config=configs/dev/node.yaml&
NODE_LISTEN=":9009" NODE_STORAGE="/tmp/node-storage-9" ./node_main run --config=configs/dev/node.yaml&

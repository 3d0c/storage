#!/bin/bash

for i in {0..9}
do
	NODE_LISTEN=":900$i" NODE_STORAGE="/tmp/node-storage-$i" build/node run --config=configs/dev/node.yaml&	
done

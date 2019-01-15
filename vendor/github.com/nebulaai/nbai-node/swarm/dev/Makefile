.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/nebulaai/nbai-node/cmd/swarm
	go build -o bin/geth     github.com/nebulaai/nbai-node/cmd/geth
	go build -o bin/bootnode github.com/nebulaai/nbai-node/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/nebulaai/nbai-node/swarm/...

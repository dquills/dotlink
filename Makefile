.PHONY: build

build:
	mkdir -p ./builds/
	go build -o builds/dotlink -ldflags "-s -w" ./cmd/dotlink/


.DEFAULT_GOAL := build
.PHONY: build

build: 
	go build -v -a -installsuffix cgo -o build/gke-api cmd/gke-api/*.go

clean:
	rm -rf build

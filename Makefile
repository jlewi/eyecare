ROOT := $(shell git rev-parse --show-toplevel)

BUILD := build

echo:
	echo ROOT=$(ROOT)

build-dir:
	mkdir -p $(BUILD)

build-go: build-dir
	CGO_ENABLED=0 go build -o $(BUILD)/eycare github.com/jlewi/eyecare/go/cmd/

tidy-go:
	gofmt -s -w .
	goimports -w .
	
lint-go:
	# golangci-lint automatically searches up the root tree for configuration files.
	golangci-lint run

test-go:
	go test -v ./...

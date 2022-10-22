ROOT := $(shell git rev-parse --show-toplevel)

BUILD := build

echo:
	echo ROOT=$(ROOT)

build-dir:
	mkdir -p $(BUILD)

build-go: build-dir
	CGO_ENABLED=0 go build -o $(BUILD)/eycare github.com/jlewi/eyecare/go/cmd/

tidy-go:
	cd go && gofmt -s -w .
	cd go && goimports -w .
	
lint-go:
	# golangci-lint automatically searches up the root tree for configuration files.
	cd go && golangci-lint run

test-go:
	cd go && go test -v ./...

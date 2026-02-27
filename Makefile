.PHONY: all build test vet fmt clean install

BINARY_NAME=todo-hours
INSTALL_PATH=$(HOME)/.local/bin

# Build flags
LDFLAGS=-ldflags "-s -w"

all: build

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/todo-hours

test:
	go test -v ./...

fmt:
	go fmt ./...

vet: fmt
	go vet ./...
	which staticcheck > /dev/null 2>&1 || go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

clean:
	rm -f $(BINARY_NAME)

install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/

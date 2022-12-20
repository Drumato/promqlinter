.PHONY: all
all: format test vet build

.PHONY: format
format:
	go fmt ./...

.PHONY: test
test:
	go test ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build -o bin/promqlinter

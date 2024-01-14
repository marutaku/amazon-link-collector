.PHONY: setup build test clean tidy


DIST_PATH = ./dist
setup:
	go mod download

build:
	go build -o ${DIST_PATH}/collector cmd/main.go

lint:
	golangci-lint run

test:
	go test -v ./...

tidy:
	go mod tidy

clean:
	rm -rf ${DIST_PATH}
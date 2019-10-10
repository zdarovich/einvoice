-include .env

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)
GOCMD=go
GOBUILD=go build
GOCLEAN=go clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=e-invoice
BINARY_DIR=bin
BINARY_UNIX=$(BINARY_NAME)_unix

all: install build

install:
	go get

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN) rm -f $(BINARY_NAME) rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./... ./$(BINARY_NAME)

# Cross compilation
build-linux:
	CGO_ENABLED=1
	GOOS=linux
	GOARCH=amd64
	$(GOBUILD) -o $(BINARY_UNIX) -v

build-windows-64:
	CGO_ENABLED=1
	GOOS=windows
	GOARCH=amd64
	$(GOBUILD) -o $(BINARY_NAME) -v

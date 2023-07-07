# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
BINARY_NAME = gget

# Check the operating system and set the binary name accordingly
ifeq ($(OS),Windows_NT)
	BINARY_NAME := $(BINARY_NAME).exe
endif

.PHONY: all build run dist clean

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v -trimpath -ldflags "-s -w" main.go

run: build
	./$(BINARY_NAME)

dist: build
	mkdir -p release
	mv $(BINARY_NAME) release/
	cp LICENSE release/

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf release

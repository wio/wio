# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=wio
BINARY_UNIX=$(BINARY_NAME)_unix

all: build run

build:
	@echo Building $(BINARY_NAME) project:
	@cd "$(CURDIR)/cmd/$(BINARY_NAME)/data" && $(GOPATH)/bin/go-bindata -nomemcopy -pkg data -prefix ../../../ ../../../assets/config/... ../../../assets/templates/...
	@cd "$(CURDIR)/cmd/$(BINARY_NAME)" && $(GOBUILD) -o $(BINARY_NAME) -v
	@if ! [ -d "bin" ]; then \
		mkdir bin; \
	fi
	@mv $(CURDIR)/cmd/${BINARY_NAME}/${BINARY_NAME} bin/
	@echo Project built!!

clean:
	@echo Cleaning $(BINARY_NAME) project files:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_UNIX)
	@echo Cleaning Finished!!

run:
	@./bin/$(BINARY_NAME) ${ARGS}

deps:
	@echo Gathering dependencies:
	$(GOGET) -t -v ./...
	@echo Dependencies up to date!

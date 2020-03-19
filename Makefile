GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BINARY_NAME=funla

all: build

race:
	go generate
	$(GOBUILD) -race -trimpath -o $(BINARY_NAME) -v .

build:
	go generate
	$(GOBUILD) -trimpath -o $(BINARY_NAME) -v .

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f stdfunfiles.go

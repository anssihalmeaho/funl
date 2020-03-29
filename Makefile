GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BINARY_NAME=funla

all: build

race:
	$(GOBUILD) -race -trimpath -o $(BINARY_NAME) -v .

build:
	$(GOBUILD) -trimpath -o $(BINARY_NAME) -v .

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=receipt-generator
VERSION=v0.0.1

all: darwin linux windows

version:
	@echo $(VERSION)
darwin: clean
	GOOS=darwin $(GOBUILD) -o release/$(BINARY_NAME)-darwin-$(VERSION) -v receipt-generator.go
linux: clean
	GOOS=linux $(GOBUILD) -o release/$(BINARY_NAME)-linux-$(VERSION) -v receipt-generator.go
windows: clean
	GOOS=windows $(GOBUILD) -o release/$(BINARY_NAME)-windows-$(VERSION).exe -v receipt-generator.go
clean:
	$(GOCLEAN)
	rm -rf release

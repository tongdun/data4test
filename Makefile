GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
GOTEST = $(GOCMD) test
BINARY_NAME = data4test
CLI = adm
TARGETDIR = ./deploy

all: serve

init:
	$(GOMOD) init $(module)

install:
	$(GOMOD) tidy

serve:
	$(GOCMD) run .

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(TARGETDIR)/$(BINARY_NAME)_linux_x86_64 main.go
    CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GOBUILD) -o $(TARGETDIR)/$(BINARY_NAME)_linux_i386 main.go
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(TARGETDIR)/$(BINARY_NAME)_darwin_x86_64 main.go
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(TARGETDIR)/$(BINARY_NAME)_linux_aarch64 main.go
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(TARGETDIR)/$(BINARY_NAME)_windows_x86_64.exe main.go
    CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(TARGETDIR)/$(BINARY_NAME)_windows_i386.exe main.go

generate:
	$(CLI) generate -c adm.ini

test: testNew black-box-test user-acceptance-test

testNew: testpre testrun testafter

testpre:
	mkdir -p biz/log
	cp config.json biz/config.json
	touch biz/log/data4test.log
	

testafter:
	rm -rf biz/log
	rm -f biz/config.json

testrun:
	$(GOTEST) -v ./...

.PHONY: all serve build generate test black-box-test user-acceptance-test ready-for-data clean
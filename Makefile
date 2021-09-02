
VERSION := $(shell git rev-parse --short HEAD)
PROJECT_NAME=listener
DEBUG=true
LOG_DIR=logs

GOPATH := $(PWD)
PKGS := $(shell ls src)

RACE_FLAG=""
ifeq ($(DEBUG),true)
	RACE_FLAG="-race"
endif

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Debug=$(DEBUG) -X=main.LogDir=$(LOG_DIR)"

install:
	@GOPATH=$(GOPATH) go get -d ./...

build:
	@GOPATH=$(GOPATH) go build $(LDFLAGS) -o bin/$(PROJECT_NAME) src/main.go

run:
	@GOPATH=$(GOPATH) go run $(RACE_FLAG) $(LDFLAGS) src/main.go

test:
	@GOPATH=$(GOPATH) go test -v

clean:
	rm -rf bin/

all:
	clean
	install
	build

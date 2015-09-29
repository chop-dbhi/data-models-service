GIT_SHA := $(shell git log -1 --pretty=format:"%h")

all: install

clean:
	go clean ./...

doc:
	godoc -http=:6060

install:
	go get github.com/julienschmidt/httprouter
	go get github.com/sirupsen/logrus
	go get github.com/russross/blackfriday
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/blang/semver

test-install: install
	go get golang.org/x/tools/cmd/cover
	go get github.com/cespare/prettybench

dev-install: install test-install

test:
	go test -cover ./...

build-assets:
	go-bindata -ignore \\.sw[a-z] -ignore \\.DS_Store assets/

build-dev-assets:
	go-bindata -debug -ignore \\.sw[a-z] -ignore \\.DS_Store assets/

_build:
	go build \
		-ldflags "-X main.progBuild=$(GIT_SHA)" \
		-o $(GOPATH)/bin/data-models .

build: build-assets _build

build-dev: build-dev-assets _build

bench:
	go test -run=none -bench=. ./... | prettybench

fmt:
	go vet ./...
	go fmt ./...

lint:
	golint ./...

.PHONY: test

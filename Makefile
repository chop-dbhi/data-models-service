GIT_SHA := $(shell git log -1 --pretty=format:"%h")

all: install

install:
	go get github.com/julienschmidt/httprouter
	go get github.com/sirupsen/logrus
	go get github.com/russross/blackfriday
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/blang/semver
	go get github.com/rs/cors

test-install: install
	go get golang.org/x/tools/cmd/cover
	go get github.com/cespare/prettybench

build-install:
	go get github.com/mitchellh/gox

test:
	go test -cover ./...

bench:
	go test -run=none -bench=. ./... | prettybench

build:
	go-bindata -debug -ignore \\.sw[a-z] -ignore \\.DS_Store assets/

	go build \
		-ldflags "-X main.progBuild=$(GIT_SHA)" \
		-o $(GOPATH)/bin/data-models .

dist-build:
	mkdir -p dist

	go-bindata -ignore \\.sw[a-z] -ignore \\.DS_Store assets/

	gox -output "dist/{{.OS}}-{{.Arch}}/data-models-service" \
		-ldflags "-X main.progBuild='$(GIT_SHA)'" \
		-os "linux windows darwin" \
		-arch "amd64" \
		. > /dev/null

dist-zip:
	cd dist && zip data-models-service-linux-amd64.zip linux-amd64/*
	cd dist && zip data-models-service-windows-amd64.zip windows-amd64/*
	cd dist && zip data-models-service-darwin-amd64.zip darwin-amd64/*


dist: dist-build dist-zip

.PHONY: test

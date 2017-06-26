ifndef GOPATH
$(error GOPATH is not set)
endif

GOARCH ?= amd64
GOOS ?= linux


clean:
	rm -rf bin/*;
	go clean -i ./...

dependencies: bootstrap
	glide install

watch: bootstrap
	goconvey

lib/lumogon_diff_plugin.so:
	mkdir -p lib
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -buildmode=plugin -o lib/lumogon_diff_plugin.so lumogon_diff_plugin.go

build: lib/lumogon_diff_plugin.so bootstrap

lint: bootstrap $(GOPATH)/src/github.com/golang/lint/golint
	golint `glide novendor`

vet: bootstrap
	go vet `glide novendor`

all: clean dependencies build

$(GOPATH)/bin/glide:
	go get -u github.com/Masterminds/glide

$(GOPATH)/src/github.com/golang/lint/golint:
	go get -u github.com/golang/lint/golint

$(GOPATH)/bin/goconvey:
	go get -u github.com/smartystreets/goconvey

bootstrap: $(GOPATH)/bin/glide $(GOPATH)/src/github.com/golang/lint/golint $(GOPATH)/bin/goconvey

.PHONY: build image test todo clean dependencies bootstrap watch

APPNAME = vasuki
VERSION=0.1.0-dev

build:
	go build -o ${APPNAME} .

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -X main.Version=${VERSION}" -v -o ${APPNAME}-linux-amd64 .

build-mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -X main.Version=${VERSION}" -v -o ${APPNAME}-darwin-amd64 .

build-all: build-mac build-linux

all: setup
	build
	install

setup:
	go get github.com/ashwanthkumar/go-gocd
	go get github.com/spf13/cobra
	go get github.com/hashicorp/go-multierror
	go get github.com/deckarep/golang-set
	go get github.com/fsouza/go-dockerclient
	go get github.com/satori/go.uuid
	go get github.com/op/go-logging
	# Test deps
	go get github.com/stretchr/testify
	go get github.com/vektra/mockery/.../

mocks:
	mockery -name=Scalar -recursive -inpkg
	mockery -name=Executor -recursive -inpkg

test:
	go test -v github.com/ashwanthkumar/vasuki/utils/sets
	go test -v github.com/ashwanthkumar/vasuki/scalar
	go test -v github.com/ashwanthkumar/vasuki

install: build
	sudo install -d /usr/local/bin
	sudo install -c ${APPNAME} /usr/local/bin/${APPNAME}

uninstall:
	sudo rm /usr/local/bin/${APPNAME}

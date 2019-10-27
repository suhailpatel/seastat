commit := $(shell git rev-parse --short=9 HEAD)
version := $(shell cat VERSION)

build:
	go build -ldflags "-X github.com/suhailpatel/seastat/flags.GitCommitHash=$(commit) -X github.com/suhailpatel/seastat/flags.Version=$(version)" -o seastat

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/suhailpatel/seastat/flags.GitCommitHash=$(commit) -X github.com/suhailpatel/seastat/flags.Version=$(version)" -o seastat-linux-$(version)

clean:
	rm -rf seastat*

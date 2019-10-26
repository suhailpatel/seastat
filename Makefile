commit := $(shell git rev-parse --short=9 HEAD)
version := $(shell cat VERSION)

build:
	go build -ldflags "-X github.com/suhailpatel/seastat/cmd.GitCommitHash=$(commit) -X github.com/suhailpatel/seastat/cmd.Version=$(version)" -o seastat

clean:
	rm -rf seastat

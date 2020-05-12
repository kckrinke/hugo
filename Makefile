#!/usr/bin/make -f

.PHONY: all clean build

help:
	@echo "usage: make <clean|build>"

clean:
	@mv -v ${GOPATH}/bin/hugo ${GOPATH}/bin/hugo.`date "+%s"`
	@go clean

build:
	@go install -v -x --tags extended

mage-fmt:
	@mage fmt

mage-build:
	@mage hugoRace

mage-test:
	@mage textCoverHTML

mage-lint:
	@mage lint
	@mage vet

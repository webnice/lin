DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

OLDGOPATH := $(GOPATH)
GOPATH := $(GOPATH)
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)
TESTPACKETS=$(shell if [ -f .testpackages ]; then cat .testpackages; fi)
BENCHPACKETS=$(shell if [ -f .benchpackages ]; then cat .benchpackages; fi)

default: lint test

link:
	@mkdir -p ${DIR}/src/gopkg.in/webnice; cd ${DIR}/src/gopkg.in/webnice && ln -s ../../.. lin.v1 2>/dev/null; true
	@if [ ! -L ${DIR}/src/vendor ]; then ln -s ${DIR}/vendor ${DIR}/src/vendor 2>/dev/null; fi
.PHONY: link

## Generate code by go generate or other utilities
generate: link
.PHONY: generate

## Dependence managers
dep: link
	@go mod download
	@go get -u
	@go mod tidy
	@go mod vendor
.PHONY: dep

test: link
	@echo "mode: set" > $(DIR)/coverage.log
	@for PACKET in $(TESTPACKETS); do \
		touch $(DIR)/coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=$(DIR)/coverage-tmp.log $$PACKET; \
		if [ "$$?" -ne "0" ]; then exit $$?; fi; \
		tail -n +2 $(DIR)/coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> $(DIR)/coverage.log; \
		rm -f $(DIR)/coverage-tmp.log; true; \
	done
.PHONY: test

cover: test
	@GOPATH=${GOPATH} go tool cover -html=$(DIR)/coverage.log
.PHONY: cover

bench: link
	@for PACKET in $(BENCHPACKETS); do GOPATH=${GOPATH} go test -race -bench=. -benchmem $$PACKET; done
.PHONY: bench

lint: link
	@gometalinter \
	--vendor \
	--deadline=15m \
	--cyclo-over=20 \
	--line-length=120 \
	--warn-unmatched-nolint \
	--disable=aligncheck \
	--enable=test \
	--enable=goimports \
	--enable=gosimple \
	--enable=misspell \
	--enable=unused \
	--enable=megacheck \
	--skip=vendor \
	--skip=src/vendor \
	--linter="vet:go tool vet -printfuncs=Infof,Debugf,Warningf,Errorf:PATH:LINE:MESSAGE" \
	./...
.PHONY: lint

clean:
	@rm -rf ${DIR}/src; true
	@rm -rf ${DIR}/bin; true
	@rm -rf ${DIR}/pkg; true
	@rm -rf ${DIR}/*.log; true
	@rm -rf ${DIR}/*.lock; true
.PHONY: clean

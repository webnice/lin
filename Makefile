DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH := $(DIR):$(GOPATH)
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)
TESTPACKETS=$(shell if [ -f .testpackages ]; then cat .testpackages; fi)
BENCHPACKETS=$(shell if [ -f .benchpackages ]; then cat .benchpackages; fi)

default: lint test

link:
	mkdir -p src/gopkg.in/webnice; cd src/gopkg.in/webnice && ln -s ../../.. nul.v1 2>/dev/null; true
.PHONY: link

## Generate code by go generate or other utilities
generate: link
.PHONY: generate

## Dependence managers
dep: link
	if command -v "gvt"; then GOPATH="$(DIR)" gvt update -all; fi
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
	GOPATH=${GOPATH} go tool cover -html=$(DIR)/coverage.log
	@make clean
.PHONY: cover

bench: link
	@for PACKET in $(BENCHPACKETS); do GOPATH=${GOPATH} go test -race -bench=. -benchmem $$PACKET; done
	@make clean
.PHONY: bench

lint: link
	gometalinter \
	--vendor \
	--deadline=15m \
	--cyclo-over=30 \
	--disable=aligncheck \
	--disable=gotype \
	--disable=structcheck \
	--skip=src/vendor \
	--linter="vet:go tool vet -printf {path}/*.go:PATH:LINE:MESSAGE" \
	./...
.PHONY: lint

clean:
	rm -rf ${DIR}/src; true
	rm -rf ${DIR}/bin/*; true
	rm -rf ${DIR}/pkg/*; true
	rm -rf ${DIR}/*.log; true
	rm -rf ${DIR}/*.lock; true
.PHONY: clean

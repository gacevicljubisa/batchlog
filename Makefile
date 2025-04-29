GO ?= go

.PHONY: binary

binary: dist FORCE
	$(GO) version
ifeq ($(OS),Windows_NT)
	$(GO) build  -o dist/batchlog.exe .
else
	$(GO) build -o dist/batchlog .
endif

dist:
	mkdir $@

FORCE:

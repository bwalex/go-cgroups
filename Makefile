GO?=go

SOURCEDIR=.

SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

.PHONY: test
test: $(SOURCES)
	$(GO) test -v -bench=.

.PHONY: clean
clean:
	rm -f $(PROG)

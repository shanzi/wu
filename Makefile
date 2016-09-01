BINARY=wu
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
VERSION=$(shell git describe --tags --always --dirty)
LD_FLAGS :=-ldflags "-X main.Version=${VERSION} -extldflags -static"

all: $(BINARY)

$(BINARY): $(SOURCES)
	go build -x --tags netgo ${LD_FLAGS} -o ${BINARY}
	strip $@
clean:
	@rm -f $(BINARY)

install:
	go install ${LD_FLAGS} -o ${BINARY}

.PHONY: clean install

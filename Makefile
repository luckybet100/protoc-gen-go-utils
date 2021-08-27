PROTO := $(shell find . -name "*.proto")

all: build install

build:
	go build

install:
	go install .

example: $(PROTO)
	protoc -I . $(PROTO) --go_out=plugins=grpc:. --go-utils_out=.

.PHONY: example clean

clean:
	rm -rf protoc-gen-go-utils
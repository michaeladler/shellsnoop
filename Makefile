PREFIX ?= /usr/local
DESTDIR ?=

BINDIR = $(DESTDIR)$(PREFIX)/bin
CC ?= clang
CFLAGS ?= -std=c99 -O2 -g -Wall -Wextra -Werror -D_FORTIFY_SOURCE=2

export CGO_ENABLED=0

.PHONY: build
build: generate shellsnoop-client
	go build -trimpath -ldflags "-s -w -X 'main.Commit=$(shell git rev-parse HEAD | tr -d [:space:])'"

.PHONY: generate
generate:
	go generate ./...

shellsnoop-client: client/main.c
	$(CC) $(CFLAGS) -o $@ $<

.PHONY: install
install:
	install -D -m0755 shellsnoop $(BINDIR)/shellsnoop
	install -D -m0755 shellsnoop-client $(BINDIR)/shellsnoop-client

.PHONY: clean
clean:
	$(RM) shellsnoop shellsnoop-client


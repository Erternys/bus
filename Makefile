GO ?= go
EXENAME = bus

GOTO_PROXY = cd ./proxy;

.PHONY: all
all: test build

%:
	mkdir $@

.PHONY: test
test:
	$(GO) test .
	$(GOTO_PROXY) cargo test

.PHONY: build
build: bin bin/bus bin/proxy

GO_BUS_FILES = $(shell find . -type f -name '*.go')
bin/bus: bin $(GO_BUS_FILES)
	$(GO) install
	$(GO) build -o ./bin/$(EXENAME) -ldflags="-s -w" ./main.go

RUST_PROXY_FILES = $(shell find proxy -type f -name '*.rs')
bin/proxy: bin $(RUST_PROXY_FILES)
	$(GOTO_PROXY)	cargo build --release; \
	cp target/release/proxy ../bin/; \

.PHONY: clean
clean:
	$(GO) clean
	$(RM) -r ./bin ./proxy/target

## (un)install script for unix

DESTDIR :=
prefix  := /usr/local
bindir  := ${prefix}/bin
mandir  := ${prefix}/share/man

.PHONY: install
install: bin/$(EXENAME)
	install -d ${DESTDIR}${bindir}
	install -m755 bin/$(EXENAME) ${DESTDIR}${bindir}/

.PHONY: uninstall
uninstall:
	rm -f ${DESTDIR}${bindir}/$(EXENAME)

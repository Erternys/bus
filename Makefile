GO ?= go
EXENAME = bus

.PHONY: all
all: build

%:
	mkdir $@

.PHONY: test
test:
	$(GO) test .

.PHONY: build
build: bin bin/proxy
	@$(GO) install
	$(GO) build -o ./bin/$(EXENAME) -ldflags="-s -w"  ./main.go

bin/proxy: bin
	cd ./proxy; \
	cargo build --release; \
	cp target/release/proxy ../bin/; \

.PHONY: clean
clean:
	$(GO) clean
	rm -rf ./bin
	rm -rf ./proxy/target

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

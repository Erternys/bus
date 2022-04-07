GO ?= go
EXENAME = bus

bin:
	mkdir bin

test:
	$(GO) test .

.PHONY: build
build: bin
	@$(GO) install
	$(GO) build -o ./bin/$(EXENAME) ./main.go

clean:
	rm -rf ./bin

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

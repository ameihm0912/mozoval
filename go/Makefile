PROJS = oval goval goval-evrtest goval-pkgquery \
	ubuntu-cve-tracker pkg2mozdef
GO = GOPATH=$(shell pwd):$(shell go env GOROOT)/bin go

all: $(PROJS)

depends:
	$(GO) get github.com/ameihm0912/gozdef
	$(GO) get github.com/ameihm0912/govfeed/src/govfeed

oval:
	$(GO) build oval
	$(GO) install oval

pkg2mozdef:
	$(GO) install pkg2mozdef

goval:
	$(GO) install goval

goval-evrtest:
	$(GO) install goval-evrtest

goval-pkgquery:
	$(GO) install goval-pkgquery

ubuntu-cve-tracker:
	$(GO) install ubuntu-cve-tracker

clean:
	rm -f bin/goval
	rm -f bin/goval-evrtest
	rm -f bin/goval-pkgquery
	rm -f pkg/*/oval.a

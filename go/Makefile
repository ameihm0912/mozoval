PROJS = oval goval goval-evrtest goval-pkgquery
GO = GOPATH=$(shell pwd):$(shell go env GOROOT)/bin go

all: $(PROJS)

oval:
	$(GO) build oval
	$(GO) install oval

goval:
	$(GO) install goval

goval-evrtest:
	$(GO) install goval-evrtest

goval-pkgquery:
	$(GO) install goval-pkgquery

clean:
	rm -f bin/goval
	rm -f pkg/*/oval.a

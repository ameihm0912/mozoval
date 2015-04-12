PROJS = oval goval goval-evrtest
GO = GOPATH=$(shell pwd):$(shell go env GOROOT)/bin go

all: $(PROJS)

oval:
	$(GO) build oval
	$(GO) install oval

goval:
	$(GO) install goval

goval-evrtest:
	$(GO) install goval-evrtest

clean:
	rm -f bin/goval
	rm -f pkg/*/oval.a

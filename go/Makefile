PROJS = oval goval goval-evrtest

all: $(PROJS)

oval:
	go build oval
	go install oval

goval:
	go install goval

goval-evrtest:
	go install goval-evrtest

clean:
	rm -f bin/goval
	rm -f pkg/*/oval.a

PROJS = oval goval

all: $(PROJS)

oval:
	go build oval
	go install oval

goval:
	go install goval

clean:
	rm -f bin/goval
	rm -f pkg/*/oval.a

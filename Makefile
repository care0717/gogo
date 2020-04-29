# Go パラメータ
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=gogo

gogo: *.go
	$(GOBUILD) -v

test: gogo
	./test.sh

clean:
	rm -f gogo *.o *~ tmp*

.PHONY: test clean

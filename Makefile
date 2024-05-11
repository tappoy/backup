PACKAGE=github.com/tappoy/backup
WORKING_DIRS=tmp bin

SRC=$(shell find . -name "*.go")
BIN=bin/$(shell basename $(CURDIR))
DOC=Document.txt
COVER=tmp/cover
COVER0=tmp/cover0

.PHONY: all clean fmt cover test lint

all: $(WORKING_DIRS) $(FMT) $(BIN) test $(DOC) lint

clean:
	rm -rf $(WORKING_DIRS)

$(WORKING_DIRS):
	mkdir -p $(WORKING_DIRS)

fmt: $(SRC)
	go fmt ./...

go.sum: go.mod
	go mod tidy

$(BIN): $(SRC) go.sum Usage.txt
	go build -o $(BIN)

test: $(BIN)
	go test -v -tags=mock -vet=all -cover -coverprofile=$(COVER)

$(DOC): $(SRC)
	go doc -all . > $(DOC)

cover: $(COVER)
	grep "0$$" $(COVER) | sed 's!$(PACKAGE)!.!' | tee $(COVER0)

lint: $(BIN)
	go vet

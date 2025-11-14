## Simple Makefile to build the m2sh binary from cmd/m2sh/main.go

BINARY := m2sh
MAIN   := cmd/m2sh/main.go

.PHONY: build clean

build: $(BINARY)

$(BINARY): $(MAIN) go.mod
	go build -o $(BINARY) $(MAIN)

clean:
	rm -f $(BINARY)

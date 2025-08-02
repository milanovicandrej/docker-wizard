
BINARY=docker-wizard
BUILD_PATH=./cmd/docker-wizard
INSTALL_PATH=/home/milanovicandrej/.local/bin/$(BINARY)
version = 0.4.0


.PHONY: all build install clean deps


all: deps build


build: deps
	go build -o $(BINARY) $(BUILD_PATH)
deps:
	go mod tidy

install: build
	sudo mv $(BINARY) $(INSTALL_PATH)

clean:
	rm -f $(BINARY)

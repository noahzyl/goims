GOCMD = go
GOBUILD  = $(GOCMD) build
BUILD_DIR = bin

BIN_SERVER=$(BUILD_DIR)/server
BIN_CLIENT=$(BUILD_DIR)/client

.PHONY: all build clean

all: build

build: $(BIN_SERVER) $(BIN_CLIENT)

$(BIN_SERVER): goims-server/main.go
	cd goims-server && $(GOBUILD) -o ../$@ main.go

$(BIN_CLIENT): goims-client/main.go
	cd goims-client && $(GOBUILD) -o ../$@ main.go

clean:
	rm -rf $(BUILD_DIR)
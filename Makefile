# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=pa
BINARY_UNIX=$(BINARY_NAME)_unix
TARGET_DIR=target

# Build flags
DEBUG_FLAGS=-gcflags="all=-N -l"
RELEASE_FLAGS=-ldflags="-s -w" -trimpath

# Build targets
all: clean build-release

build-debug: prepare
	$(GOBUILD) $(DEBUG_FLAGS) -o $(TARGET_DIR)/debug/$(BINARY_NAME)-checker ./checker
	$(GOBUILD) $(DEBUG_FLAGS) -o $(TARGET_DIR)/debug/$(BINARY_NAME)-batch-checker ./batch_checker
	$(GOBUILD) $(DEBUG_FLAGS) -o $(TARGET_DIR)/debug/$(BINARY_NAME)-encoder ./encoder
	$(GOBUILD) $(DEBUG_FLAGS) -o $(TARGET_DIR)/debug/$(BINARY_NAME)-evmaddress ./evmaddress_generator
	$(GOBUILD) $(DEBUG_FLAGS) -o $(TARGET_DIR)/debug/$(BINARY_NAME)-server ./server

build-release: prepare
	$(GOBUILD) $(RELEASE_FLAGS) -o $(TARGET_DIR)/release/$(BINARY_NAME)-checker ./checker
	$(GOBUILD) $(RELEASE_FLAGS) -o $(TARGET_DIR)/release/$(BINARY_NAME)-batch-checker ./batch_checker
	$(GOBUILD) $(RELEASE_FLAGS) -o $(TARGET_DIR)/release/$(BINARY_NAME)-encoder ./encoder
	$(GOBUILD) $(RELEASE_FLAGS) -o $(TARGET_DIR)/release/$(BINARY_NAME)-evmaddress ./evmaddress_generator
	$(GOBUILD) $(RELEASE_FLAGS) -o $(TARGET_DIR)/release/$(BINARY_NAME)-server ./server

prepare:
	mkdir -p $(TARGET_DIR)/debug
	mkdir -p $(TARGET_DIR)/release

clean:
	$(GOCLEAN)
	rm -rf $(TARGET_DIR)

# Cross compilation
build-linux-debug: prepare
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(DEBUG_FLAGS) -o $(TARGET_DIR)/debug/$(BINARY_UNIX)-checker ./checker
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(DEBUG_FLAGS) -o $(TARGET_DIR)/debug/$(BINARY_UNIX)-batch-checker ./batch_checker
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(DEBUG_FLAGS) -o $(TARGET_DIR)/debug/$(BINARY_UNIX)-encoder ./encoder
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(DEBUG_FLAGS) -o $(TARGET_DIR)/debug/$(BINARY_UNIX)-evmaddress ./evmaddress_generator
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(DEBUG_FLAGS) -o $(TARGET_DIR)/debug/$(BINARY_UNIX)-server ./server

build-linux-release: prepare
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(RELEASE_FLAGS) -o $(TARGET_DIR)/release/$(BINARY_UNIX)-checker ./checker
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(RELEASE_FLAGS) -o $(TARGET_DIR)/release/$(BINARY_UNIX)-batch-checker ./batch_checker
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(RELEASE_FLAGS) -o $(TARGET_DIR)/release/$(BINARY_UNIX)-encoder ./encoder
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(RELEASE_FLAGS) -o $(TARGET_DIR)/release/$(BINARY_UNIX)-evmaddress ./evmaddress_generator
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(RELEASE_FLAGS) -o $(TARGET_DIR)/release/$(BINARY_UNIX)-server ./server

.PHONY: all build-debug build-release clean build-linux-debug build-linux-release prepare
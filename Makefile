# Name of the binary
BINARY=metronome-cli
BUILD_DIR=bin

# Default build
all: build

# Build binary
build:
	@mkdir -p $(BUILD_DIR)
	go mod tidy
	go build -o $(BUILD_DIR)/$(BINARY) .

# Install to /usr/local/bin
install: build
	cp $(BUILD_DIR)/$(BINARY) /usr/local/bin/$(BINARY)
	chmod +x /usr/local/bin/$(BINARY)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Run the metronome
run: build
	./$(BUILD_DIR)/$(BINARY)

.PHONY: all build install clean run
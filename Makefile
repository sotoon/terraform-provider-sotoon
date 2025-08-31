PROVIDER_SOURCE := sotoon/sotoon
PROVIDER_BINARY_NAME := terraform-provider-sotoon
INSTALL_DIR := $(shell go env GOPATH)/bin
OS_ARCH := $(shell go env GOOS)_$(shell go env GOARCH)


.PHONY: all build install configure-tf-rc clean help refresh

all: build install

# Build the provider binary
build:
	@echo "--> Building $(PROVIDER_BINARY_NAME)..."
	go build -o $(PROVIDER_BINARY_NAME)

# Install the built provider binary to the specified INSTALL_DIR
install: build
	@echo "--> Installing $(PROVIDER_BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@mv $(PROVIDER_BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(PROVIDER_BINARY_NAME)
	@echo "    Installed: $(INSTALL_DIR)/$(PROVIDER_BINARY_NAME)"

# NEW: Clean and rebuild the provider
refresh: clean install
	@echo "--> Provider refreshed successfully."

# Clean up compiled binary and optionally the installed one
clean:
	@echo "--> Cleaning up..."
	@rm -f $(PROVIDER_BINARY_NAME) # Removes the binary from the build directory
	@if [ -f "$(INSTALL_DIR)/$(PROVIDER_BINARY_NAME)" ]; then \
	  echo "    Removing installed binary from $(INSTALL_DIR)..."; \
	  rm $(INSTALL_DIR)/$(PROVIDER_BINARY_NAME); \
	else \
	  echo "    Installed binary not found at $(INSTALL_DIR), skipping removal."; \
	fi
	@echo "    Remember to manually clean up your .terraformrc if needed."

help:
	@echo "Usage:"
	@echo "  make refresh            - Cleans, rebuilds, and installs the provider"
	@echo "  make all                - Builds, installs, and configures .terraformrc"
	@echo "  make build              - Only builds the provider binary"
	@echo "  make install            - Installs the built provider to $(INSTALL_DIR)"
	@echo "  make clean              - Removes compiled binary from current dir and $(INSTALL_DIR)"
	@echo "  make help               - Displays this help message"
# Makefile extensions for darwin.

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

SENZING_DIR ?= /opt/senzing/g2
SENZING_TOOLS_SENZING_DIRECTORY ?= $(SENZING_DIR)

LD_LIBRARY_PATH := $(SENZING_TOOLS_SENZING_DIRECTORY)/lib:$(SENZING_TOOLS_SENZING_DIRECTORY)/lib/macos
DYLD_LIBRARY_PATH := $(LD_LIBRARY_PATH)

# -----------------------------------------------------------------------------
# OS specific targets
# -----------------------------------------------------------------------------

.PHONY: clean-osarch-specific
clean-osarch-specific:
	@rm -rf $(TARGET_DIRECTORY) || true
	@rm -f $(GOPATH)/bin/$(PROGRAM_NAME) || true


.PHONY: hello-world-osarch-specific
hello-world-osarch-specific:
	@echo "Hello World, from darwin."


.PHONY: run-osarch-specific
run-osarch-specific:
	@go run -exec macos_exec_dyld.sh main.go


.PHONY: setup-osarch-specific
setup-osarch-specific:
	@echo "No setup required."


.PHONY: test-osarch-specific
test-osarch-specific:
	@go test -exec macos_exec_dyld.sh -v -p 1 ./...

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-darwin
only-darwin:
	@echo "Only darwin has this Makefile target."
# Makefile extensions for windows.

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------


# -----------------------------------------------------------------------------
# OS specific targets
# -----------------------------------------------------------------------------

.PHONY: clean-osarch-specific
clean-osarch-specific:
	@del /F /S /Q $(GOPATH)/bin/$(PROGRAM_NAME)
	@del /F /S /Q $(MAKEFILE_DIRECTORY)/coverage.html
	@del /F /S /Q $(MAKEFILE_DIRECTORY)/coverage.out
	@del /F /S /Q $(MAKEFILE_DIRECTORY)/cover.out
	@del /F /S /Q $(TARGET_DIRECTORY)
	@taskkill /f /t/im godoc


.PHONY: coverage-osarch-specific
coverage-osarch-specific:
	@go test -v -coverprofile=coverage.out -p 1 ./...
	@go tool cover -html="coverage.out" -o coverage.html
	@explorer file://$(MAKEFILE_DIRECTORY)/coverage.html


.PHONY: documentation-osarch-specific
documentation-osarch-specific:
	@taskkill /f /t/im godoc
	@start /b godoc
	@explorer http://localhost:6060


.PHONY: hello-world-osarch-specific
hello-world-osarch-specific:
	$(info Hello World, from windows.)


.PHONY: run-osarch-specific
run-osarch-specific:
	@go run main.go


.PHONY: setup-osarch-specific
setup-osarch-specific:
	$(info No setup required.)


.PHONY: test-osarch-specific
test-osarch-specific:
	@go test -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-windows
only-windows:
	$(info Only windows has this Makefile target.)

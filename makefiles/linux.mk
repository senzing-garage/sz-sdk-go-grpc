# Makefile extensions for linux.

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------


# -----------------------------------------------------------------------------
# OS specific targets
# -----------------------------------------------------------------------------

.PHONY: clean-osarch-specific
clean-osarch-specific:
	@docker rm --force senzing-tools-serve-grpc || true
	@rm -rf $(TARGET_DIRECTORY) || true
	@rm -f $(GOPATH)/bin/$(PROGRAM_NAME) || true


.PHONY: hello-world-osarch-specific
hello-world-osarch-specific:
	@echo "Hello World, from linux."


.PHONY: run-osarch-specific
run-osarch-specific:
	@go run main.go


.PHONY: setup-osarch-specific
setup-osarch-specific:
	@docker run \
		--detach \
		--env SENZING_TOOLS_COMMAND=serve-grpc \
		--env SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db \
		--env SENZING_TOOLS_ENABLE_ALL=true \
		--name senzing-tools-serve-grpc \
		--publish 8261:8261 \
		--rm \
		senzing/senzing-tools
	@echo "senzing/senzing-tools server-grpc running in background."


.PHONY: test-osarch-specific
test-osarch-specific:
	@go test -v -p 1 ./...

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-linux
only-linux:
	@echo "Only linux has this Makefile target."

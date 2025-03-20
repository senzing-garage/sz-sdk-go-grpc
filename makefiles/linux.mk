# Makefile extensions for linux.

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

PATH := $(MAKEFILE_DIRECTORY)/bin:/$(HOME)/go/bin:$(PATH)

# -----------------------------------------------------------------------------
# OS specific targets
# -----------------------------------------------------------------------------

.PHONY: clean-osarch-specific
clean-osarch-specific:
	@docker rm --force senzing-serve-grpc || true
	@rm -f  $(GOPATH)/bin/$(PROGRAM_NAME) || true
	@rm -f  $(MAKEFILE_DIRECTORY)/.coverage || true
	@rm -f  $(MAKEFILE_DIRECTORY)/coverage.html || true
	@rm -f  $(MAKEFILE_DIRECTORY)/coverage.out || true
	@rm -f  $(MAKEFILE_DIRECTORY)/cover.out || true
	@rm -fr $(TARGET_DIRECTORY) || true
	@rm -fr /tmp/sqlite || true
	@pkill godoc || true


.PHONY: coverage-osarch-specific
coverage-osarch-specific: export SENZING_LOG_LEVEL=TRACE
coverage-osarch-specific:
	@go test -v -coverprofile=coverage.out -p 1 ./...
	@go tool cover -html="coverage.out" -o coverage.html
	@xdg-open $(MAKEFILE_DIRECTORY)/coverage.html


.PHONY: dependencies-for-development-osarch-specific
dependencies-for-development-osarch-specific:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest


.PHONY: documentation-osarch-specific
documentation-osarch-specific:
	@pkill godoc || true
	@godoc &
	@xdg-open http://localhost:6060


.PHONY: hello-world-osarch-specific
hello-world-osarch-specific:
	$(info Hello World, from linux.)


.PHONY: run-osarch-specific
run-osarch-specific:
	@go run main.go


.PHONY: setup-osarch-specific
setup-osarch-specific:
	@docker run \
		--detach \
		--env SENZING_TOOLS_ENABLE_ALL=true \
		--name senzing-serve-grpc \
		--publish 8261:8261 \
		--rm \
		senzing/serve-grpc
	$(info senzing/serve-grpc running in background.)


.PHONY: setup-mutual-tls-osarch-specific
setup-mutual-tls-osarch-specific:
	@docker run \
		--detach \
		--env SENZING_TOOLS_CLIENT_CA_CERTIFICATE_FILE=/testdata/certificates/certificate-authority/certificate.pem \
		--env SENZING_TOOLS_ENABLE_ALL=true \
		--env SENZING_TOOLS_SERVER_CERTIFICATE_FILE=/testdata/certificates/server/certificate.pem \
		--env SENZING_TOOLS_SERVER_KEY_FILE=/testdata/certificates/server/private_key.pem \
		--name senzing-serve-grpc \
		--publish 8261:8261 \
		--rm \
		--volume $(MAKEFILE_DIRECTORY)/testdata:/testdata \
		senzing/serve-grpc
	$(info senzing/serve-grpc with Mutual TLS running in background.)


.PHONY: setup-server-side-tls-osarch-specific
setup-server-side-tls-osarch-specific:
	@docker run \
		--detach \
		--env SENZING_TOOLS_ENABLE_ALL=true \
		--env SENZING_TOOLS_SERVER_CERTIFICATE_FILE=/testdata/certificates/server/certificate.pem \
		--env SENZING_TOOLS_SERVER_KEY_FILE=/testdata/certificates/server/private_key.pem \
		--name senzing-serve-grpc \
		--publish 8261:8261 \
		--rm \
		--volume $(MAKEFILE_DIRECTORY)/testdata:/testdata \
		senzing/serve-grpc
	$(info senzing/serve-grpc with Server-Side TLS running in background.)


.PHONY: test-osarch-specific
test-osarch-specific:
	@go test -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: test-mutual-tls-osarch-specific
test-mutual-tls-osarch-specific: export SENZING_TOOLS_SERVER_CA_CERTIFICATE_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/certificate-authority/certificate.pem
test-mutual-tls-osarch-specific: export SENZING_TOOLS_CLIENT_CERTIFICATE_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/client/certificate.pem
test-mutual-tls-osarch-specific: export SENZING_TOOLS_CLIENT_KEY_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/client/private_key.pem
test-mutual-tls-osarch-specific:
	@go test -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: test-mutual-tls-encrypted-key-osarch-specific
test-mutual-tls-encrypted-key-osarch-specific: export SENZING_TOOLS_SERVER_CA_CERTIFICATE_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/certificate-authority/certificate.pem
test-mutual-tls-encrypted-key-osarch-specific: export SENZING_TOOLS_CLIENT_CERTIFICATE_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/client/certificate.pem
test-mutual-tls-encrypted-key-osarch-specific: export SENZING_TOOLS_CLIENT_KEY_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/client/private_key_encrypted.pem
test-mutual-tls-encrypted-key-osarch-specific: export SENZING_TOOLS_CLIENT_KEY_PASSPHRASE=Passw0rd
test-mutual-tls-encrypted-key-osarch-specific:
	@go test -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: test-server-side-tls-osarch-specific
test-server-side-tls-osarch-specific: export SENZING_TOOLS_SERVER_CA_CERTIFICATE_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/certificate-authority/certificate.pem
test-server-side-tls-osarch-specific:
	@go test -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: test-mutual-tls
test-mutual-tls: test-mutual-tls-osarch-specific

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-linux
only-linux:
	$(info Only linux has this Makefile target.)

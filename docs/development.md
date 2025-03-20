# sz-sdk-go-grpc development

The following instructions are useful during development.

**Note:** This has been tested on Linux and Darwin/macOS.
It has not been tested on Windows.

## Prerequisites for development

:thinking: The following tasks need to be complete before proceeding.
These are "one-time tasks" which may already have been completed.

1. The following software programs need to be installed:
    1. [git]
    1. [make]
    1. [docker]
    1. [go]

## Install Git repository

1. Identify git repository.

    ```console
    export GIT_ACCOUNT=senzing-garage
    export GIT_REPOSITORY=sz-sdk-go-grpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Using the environment variables values just set, follow
   steps in [clone-repository] to install the Git repository.

## Dependencies

1. A one-time command to install dependencies needed for `make` targets.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make dependencies-for-development

    ```

1. Install dependencies needed for [Go] code.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make dependencies

    ```

## Lint

1. Run linting.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make lint

    ```

## Test

1. Run tests.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean setup test

    ```

### Test Server-Side TLS

1. Run a gRPC server.
   Either:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean setup-server-side-tls test-server-side-tls
    ```

## Test Mutual TLS

1. Run a gRPC server.
   Either:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean setup-mutual-tls test-mutual-tls
    ```

   Or:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean setup-mutual-tls test-mutual-tls-encrypted-key
    ```

## Coverage

Create a code coverage map.

1. Run Go tests.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean setup coverage

    ```

   A web-browser will show the results of the coverage.
   The goal is to have over 80% coverage.
   Anything less needs to be reflected in [testcoverage.yaml].

## Documentation

1. View documentation.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean documentation

    ```

1. If a web page doesn't appear, visit [localhost:6060].
1. Senzing documentation will be in the "Third party" section.
   `github.com` > `senzing-garage` > `sz-sdk-go-grpc`

1. When a versioned release is published with a `v0.0.0` format tag,
the reference can be found by clicking on the following badge at the top of the README.md page.
Example:

    [![Go Reference Badge]][Go Reference]

1. To stop the `godoc` server, run

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean

    ```

## References

[clone-repository]: https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/clone-repository.md
[docker]: https://github.com/senzing-garage/knowledge-base/blob/main/WHATIS/docker.md
[git]: https://github.com/senzing-garage/knowledge-base/blob/main/WHATIS/git.md
[Go Reference Badge]: https://pkg.go.dev/badge/github.com/senzing-garage/sz-sdk-go-grpc.svg
[Go Reference]: https://pkg.go.dev/github.com/senzing-garage/sz-sdk-go-grpc
[go]: https://github.com/senzing-garage/knowledge-base/blob/main/WHATIS/go.md
[localhost:6060]: http://localhost:6060/pkg/github.com/senzing-garage/sz-sdk-go-grpc/
[make]: https://github.com/senzing-garage/knowledge-base/blob/main/WHATIS/make.md
[testcoverage.yaml]: ../.github/coverage/testcoverage.yaml

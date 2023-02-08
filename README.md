# g2-sdk-go-grpc

## :warning: WARNING: g2-sdk-go-grpc is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing g2-sdk-go-grpc packages provide a Software Development Kit that
communicates with a
[Senzing gRPC server](https://github.com/Senzing/servegrpc).
`g2-sdk-go-grpc` is one of the implementations returned by the
[Senzing/go-sdk-abstract-factory](https://github.com/Senzing/go-sdk-abstract-factory)

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/g2-sdk-go-grpc.svg)](https://pkg.go.dev/github.com/senzing/g2-sdk-go-grpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/g2-sdk-go-grpc)](https://goreportcard.com/report/github.com/senzing/g2-sdk-go-grpc)
[![go-test.yaml](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/go-test.yaml/badge.svg)](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/go-test.yaml)

## Overview

The Senzing g2-sdk-go-grpc packages enable Go programs to call Senzing library functions
across a network to a
[Senzing gRPC server](https://github.com/Senzing/servegrpc).

Just like
[g2-sdk-go](https://github.com/Senzing/g2-sdk-go),
the `g2-sdk-go-grpc` packages implement the following interfaces:

1. [G2config](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2config#G2config)
1. [G2configmgr](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2configmgr#G2configmgr)
1. [G2diagnostic](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2diagnostic#G2diagnostic)
1. [G2engine](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2engine#G2engine)
1. [G2product](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2product#G2product)

## Development

### Run a Senzing gRPC server

To run a Senzing gRPC server, visit
[Senzing/servegrpc](https://github.com/Senzing/servegrpc).

### Build

The following instructions build the example `main.go` program.

1. Identify repository.
   Example:

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=g2-sdk-go-grpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Using the environment variables values just set, follow steps in [clone-repository](https://github.com/Senzing/knowledge-base/blob/main/HOWTO/clone-repository.md) to install the Git repository.
1. Build the binaries.
   Example:

     ```console
     cd ${GIT_REPOSITORY_DIR}
     make build

     ```

1. The binaries will be found in ${GIT_REPOSITORY_DIR}/target.
   Example:

    ```console
    tree ${GIT_REPOSITORY_DIR}/target

    ```

1. Run the binary.
   Example:

    ```console
    ${GIT_REPOSITORY_DIR}/target/linux/template-go

    ```

1. Clean up.
   Example:

     ```console
     cd ${GIT_REPOSITORY_DIR}
     make clean

     ```

### Test

1. Identify git repository.

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=g2-sdk-go-grpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Using the environment variables values just set, follow steps in
   [clone-repository](https://github.com/Senzing/knowledge-base/blob/main/HOWTO/clone-repository.md) to install the Git repository.

1. Run tests.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make test

    ```

### Run test cases

1. Identify git repository.

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=g2-sdk-go
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Using the environment variables values just set, follow steps in
   [clone-repository](https://github.com/Senzing/knowledge-base/blob/main/HOWTO/clone-repository.md) to install the Git repository.

1. Set environment variables.
   Identify Database URL of database in docker-compose stack.
   Example:

    ```console
    export LOCAL_IP_ADDRESS=$(curl --silent https://raw.githubusercontent.com/Senzing/knowledge-base/main/gists/find-local-ip-address/find-local-ip-address.py | python3 -)
    export SENZING_TOOLS_DATABASE_URL=postgresql://postgres:postgres@${LOCAL_IP_ADDRESS}:5432/G2

    ```

1. Run tests.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make test

    ```

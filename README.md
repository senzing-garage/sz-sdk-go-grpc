# g2-sdk-go-grpc

## :warning: WARNING: g2-sdk-go-grpc is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing `g2-sdk-go-grpc` packages provide a Go Software Development Kit
adhering to the
[g2-sdk-go](https://github.com/Senzing/g2-sdk-go) interfaces that
communicates with a
[Senzing gRPC server](https://github.com/Senzing/servegrpc).

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/g2-sdk-go-grpc.svg)](https://pkg.go.dev/github.com/senzing/g2-sdk-go-grpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/g2-sdk-go-grpc)](https://goreportcard.com/report/github.com/senzing/g2-sdk-go-grpc)
[![go-test.yaml](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/go-test.yaml/badge.svg)](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/go-test.yaml)

## Overview

The Senzing `g2-sdk-go-grpc` packages enable Go programs to call Senzing library functions
across a network to a
[Senzing gRPC server](https://github.com/Senzing/servegrpc).

`g2-sdk-go-grpc` packages implement the following
[g2-sdk-go](https://github.com/Senzing/g2-sdk-go)
interfaces:

1. [G2config](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2api#G2config)
1. [G2configmgr](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2api#G2configmgr)
1. [G2diagnostic](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2api#G2diagnostic)
1. [G2engine](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2api#G2engine)
1. [G2product](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2api#G2product)

Other implementations of the
[g2-sdk-go](https://github.com/Senzing/g2-sdk-go)
interface include:

- [g2-sdk-go-base](https://github.com/Senzing/g2-sdk-go-base) - for
  calling Senzing SDK APIs natively
- [g2-sdk-go-mock](https://github.com/Senzing/g2-sdk-go-mock) - for
  unit testing calls to the Senzing Go SDK
- [go-sdk-abstract-factory](https://github.com/Senzing/go-sdk-abstract-factory) - An
  [abstract factory pattern](https://en.wikipedia.org/wiki/Abstract_factory_pattern)
  for switching among implementations

## Use

(TODO:)

## Development

### Install Git repository

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

### Run a Senzing gRPC server

To run a Senzing gRPC server, visit
[Senzing/servegrpc](https://github.com/Senzing/servegrpc).

### Test

1. Run tests.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean test

    ```

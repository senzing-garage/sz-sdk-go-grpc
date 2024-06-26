# sz-sdk-go-grpc

If you are beginning your journey with [Senzing],
please start with [Senzing Quick Start guides].

You are in the [Senzing Garage] where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: sz-sdk-go-grpc is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing `sz-sdk-go-grpc` packages provide a [Go]
language Software Development Kit adhering to the
[sz-sdk-go] interfaces that communicates with a [Senzing gRPC server].

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing-garage/g2-sdk-go-grpc.svg)](https://pkg.go.dev/github.com/senzing-garage/g2-sdk-go-grpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing-garage/g2-sdk-go-grpc)](https://goreportcard.com/report/github.com/senzing-garage/g2-sdk-go-grpc)
[![License](https://img.shields.io/badge/License-Apache2-brightgreen.svg)](https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/LICENSE)

[![go-test-linux.yaml](https://github.com/senzing-garage/g2-sdk-go-grpc/actions/workflows/go-test-linux.yaml/badge.svg)](https://github.com/senzing-garage/g2-sdk-go-grpc/actions/workflows/go-test-linux.yaml)
[![go-test-darwin.yaml](https://github.com/senzing-garage/g2-sdk-go-grpc/actions/workflows/go-test-darwin.yaml/badge.svg)](https://github.com/senzing-garage/g2-sdk-go-grpc/actions/workflows/go-test-darwin.yaml)
[![go-test-windows.yaml](https://github.com/senzing-garage/g2-sdk-go-grpc/actions/workflows/go-test-windows.yaml/badge.svg)](https://github.com/senzing-garage/g2-sdk-go-grpc/actions/workflows/go-test-windows.yaml)

## Overview

The Senzing `g2-sdk-go-grpc` packages enable Go programs to call Senzing library functions
across a network to a
[Senzing gRPC server](https://github.com/senzing-garage/servegrpc).

Other implementations of the [sz-sdk-go]
interface include:

- [sz-sdk-go-core] - for calling Senzing SDK APIs natively
- [sz-sdk-go-mock] - for unit testing calls to the Senzing Go SDK
- [go-sdk-abstract-factory] - An [abstract factory pattern] for switching among implementations

## Use

(TODO:)

## References

1. [Development]
1. [Errors]
1. [Examples]
1. [Package reference]

[abstract factory pattern]: https://en.wikipedia.org/wiki/Abstract_factory_pattern
[Development]: docs/development.md
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[go-sdk-abstract-factory]: https://github.com/senzing-garage/go-sdk-abstract-factory
[Go]: https://go.dev/
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/sz-sdk-go-core
[Senzing Garage]: https://github.com/senzing-garage-garage
[Senzing gRPC server]: https://github.com/senzing-garage/servegrpc
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[Senzing]: https://senzing.com/
[sz-sdk-go-core]: https://github.com/senzing-garage/sz-sdk-go-core
[sz-sdk-go-mock]: https://github.com/senzing-garage/sz-sdk-go-mock
[sz-sdk-go]: https://github.com/senzing-garage/sz-sdk-go

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
[markdownlint](https://dlaa.me/markdownlint/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

-

## [0.3.3] - 2023-09-01

### Changed in 0.3.3

- Last version before SenzingAPI 3.8.0

## [0.3.2] - 2023-08-05

### Changed in 0.3.2

- Changed default port to 8261
- Moved to `go-logging`
- Refactor to `template-go`
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.8
  - github.com/senzing/go-common v0.2.11
  - github.com/senzing/go-logging v1.3.2
  - github.com/senzing/go-observing v0.2.7
  - google.golang.org/grpc v1.57.0

## [0.3.1] - 2023-05-26

### Changed in 0.3.1

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.5
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20230608182106-25c8cdc02e3c
  - github.com/senzing/go-common v0.1.4
  - github.com/senzing/go-logging v1.2.6
  - github.com/senzing/go-observing v0.2.6
  - github.com/stretchr/testify v1.8.4
  - google.golang.org/grpc v1.56.0

## [0.3.0] - 2023-05-26

### Changed in 0.3.0

- Change `g2config.Load()` signature
- Added `gosec`
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.4
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20230526140633-b44eb0f20e1b

## [0.2.7] - 2023-05-19

### Changed in 0.2.7

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.3

## [0.2.6] - 2023-05-11

### Changed in 0.2.6

- Update dependencies
  - github.com/senzing/go-common v0.1.3
  - github.com/senzing/go-logging v1.2.3

## [0.2.5] - 2023-05-10

### Changed in 0.2.5

- Added `GetObserverOrigin()` and `SetObserverOrigin()` to g2* packages
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.2
  - github.com/senzing/go-observing v0.2.2
  - google.golang.org/grpc v1.55.0

## [0.2.4] - 2023-04-22

### Changed in 0.2.4

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.1

## [0.2.3] - 2023-04-21

### Changed in 0.2.3

- Changed `SetLogLevel(ctx context.Context, logLevel logger.Level)` to `SetLogLevel(ctx context.Context, logLevelName string)`

## [0.2.2] - 2023-03-29

### Changed in 0.2.2

- Added `helper.ConvertGrpcError()`
- Refactored documentation
- Updated dependencies

## [0.2.1] - 2023-02-21

### Changed in 0.2.1

- Change GetSdkId() signature.

## [0.2.0] - 2023-02-15

### Changed in 0.2.0

- Using refactored g2-sdk-go

## [0.1.0] - 2023-01-31

### Added to 0.1.0

- Initial implementation

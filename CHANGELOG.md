# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
[markdownlint](https://dlaa.me/markdownlint/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

-

## [0.2.5] - 2023-05-10

### Changed in 0.2.5

- Added GetObserverOrigin() and SetObserverOrigin() to g2* packages
- Update dependencies
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

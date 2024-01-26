# g2-sdk-go-grpc development

## Install Go

1. See Go's [Download and install](https://go.dev/doc/install)

## Run a Senzing gRPC server

To run a Senzing gRPC server, visit
[Senzing/servegrpc](https://github.com/senzing-garage/servegrpc).

A simple method using `senzing-tools`.

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
export SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db
senzing-tools init-database
senzing-tools serve-grpc

```

A simple method using `senzing/senzing-tools` Docker image.

```console
docker run \
    --env SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db \
    --publish 8261:8261 \
    --rm \
    senzing/senzing-tools serve-grpc --enable-all

```

## Install Git repository

The following instructions build the example `main.go` program.

1. Identify repository.
   Example:

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=g2-sdk-go-grpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Using the environment variables values just set, follow steps in [clone-repository](https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/clone-repository.md) to install the Git repository.

## Test

1. Run tests.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean test

    ```
